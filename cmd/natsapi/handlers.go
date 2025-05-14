package natsapi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
)

// subscriptionHandler обработчик подписки приёма запросов
func (api *apiNatsModule) subscriptionRequestHandler() {
	_, err := api.natsConn.Subscribe(api.subscriptionRequest, func(m *nats.Msg) {

		/*

			написать хранилище запросов, где ключ - идентификатор задачи
			значение - все искомые адреса

		*/

		api.chFromModule <- SettingsChanOutput{
			TaskId: uuid.NewString(),
			Data:   m.Data,
		}

		//счетчик принятых запросов
		api.counter.SendMessage("update accepted events", 1)
	})
	if err != nil {
		api.logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
}

// incomingInformationHandler обработчик информации полученной изнутри приложения
func (api *apiNatsModule) incomingInformationHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case incomingData := <-api.chToModule:
			//команда на установку тега
			if err := api.natsConn.Publish(api.settings.command,
				fmt.Appendf(nil, `{
									  "service": "placeholder_docbase_db",
									  "command": "add_case_tag",
									  "root_id": "%s",
									  "case_id": "%s",
									  "value": "Webhook: send=\"ElasticsearchDB"
								}`, incomingData.RootId, incomingData.CaseId)); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}
		}
	}
}
