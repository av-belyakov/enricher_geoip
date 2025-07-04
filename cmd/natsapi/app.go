// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_geoip/cmd/natsapi/storagetemporary"
	"github.com/av-belyakov/enricher_geoip/interfaces"
	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(counter interfaces.Counter, logger interfaces.Logger, opts ...NatsApiOptions) (*apiNatsModule, error) {
	api := &apiNatsModule{
		settings: apiNatsSettings{
			cachettl: 10,
		},
		//для подсчёта
		counter: counter,
		//для логирования
		logger: logger,
		//запросы в модуль
		chFromModule: make(chan interfaces.Requester),
		//события из модуля
		chToModule: make(chan interfaces.Responser),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// Start инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiNatsModule) Start(ctx context.Context) error {
	storage, err := storagetemporary.New(
		ctx,
		storagetemporary.WithCacheTTL(300),
		storagetemporary.WithCacheTimeTick(10),
	)
	if err != nil {
		return err
	}

	api.storage = storage

	if ctx.Err() != nil {
		return ctx.Err()
	}

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.settings.host, api.settings.port),
		//имя клиента
		nats.Name(fmt.Sprintf("enricher_geoip.%s", api.settings.nameRegionalObject)),
		//неограниченное количество попыток переподключения
		nats.MaxReconnects(-1),
		//время ожидания до следующей попытки переподключения (по умолчанию 2 сек.)
		nats.ReconnectWait(3*time.Second),
		//обработка разрыва соединения с NATS
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected (%w)", err)).Error())
		}),
		//обработка переподключения к NATS
		nats.ReconnectHandler(func(c *nats.Conn) {
			api.logger.Send("info", "the connection to NATS has been re-established")
		}),
		//поиск медленных получателей (не обязательный для данного приложения параметр)
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			if err == nats.ErrSlowConsumer {
				pendingMsgs, _, err := s.Pending()
				if err != nil {
					api.logger.Send("warning", fmt.Sprintf("couldn't get pending messages: %v", err))

					return
				}

				api.logger.Send("warning", fmt.Sprintf("Falling behind with %d pending messages on subject %q.\n", pendingMsgs, s.Subject))
			}
		}))
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	api.natsConn = nc

	//обработчик подписки запросов поиска информации по БД GeoIP
	go api.subscriptionRequestHandler()

	//обработчик информации полученной изнутри приложения
	go api.incomingInformationHandler(ctx)

	context.AfterFunc(ctx, func() {
		api.storage.Cancel()
		nc.Drain()
	})

	return nil
}
