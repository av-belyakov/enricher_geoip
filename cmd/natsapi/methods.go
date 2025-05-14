package natsapi

import "errors"

// GetChanDataToModule канал для передачи данных в модуль
func (api *apiNatsModule) GetChanDataToModule() chan SettingsChanInput {
	return api.chToModule
}

// GetChanDataFromModule канал для приёма данных из модуля
func (api *apiNatsModule) GetChanDataFromModule() chan SettingsChanOutput {
	return api.chFromModule
}

//******************* функции настройки опций natsapi ***********************

// WithHost имя или ip адрес хоста API
func WithHost(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.settings.host = v

		return nil
	}
}

// WithPort порт API
func WithPort(v int) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.settings.port = v

		return nil
	}
}

// WithCacheTTL время жизни для кэша хранящего функции-обработчики запросов к модулю
func WithCacheTTL(v int) NatsApiOptions {
	return func(th *apiNatsModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.settings.cachettl = v

		return nil
	}
}

// WithNameRegionalObject наименование которое будет отображатся в статистике подключений NATS
func WithNameRegionalObject(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		n.settings.nameRegionalObject = v

		return nil
	}
}

// WithSubscriptionRequest 'слушатель' запросов на поиск информации
func WithSubscriptionRequest(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'subscription_request' cannot be empty")
		}

		n.subscriptionRequest = v

		return nil
	}
}

// WithSubscriptionResponse подписка для передачи ответов
func WithSubscriptionResponse(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'subscription_response' cannot be empty")
		}

		n.subscriptionResponse = v

		return nil
	}
}

// WithSendCommand команду отправляемая в NATS
func WithSendCommand(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'command' cannot be empty")
		}

		n.settings.command = v

		return nil
	}
}
