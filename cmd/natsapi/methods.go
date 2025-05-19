package natsapi

import (
	"errors"

	"github.com/av-belyakov/enricher_geoip/interfaces"
)

// GetChToModule канал для передачи данных в модуль
func (api *apiNatsModule) GetChToModule() chan interfaces.Responser {
	return api.chToModule
}

// GetChFromModule канал для приёма данных из модуля
func (api *apiNatsModule) GetChFromModule() chan interfaces.Requester {
	return api.chFromModule
}

//--- для ObjectFromNats ---

func (o *ObjectFromNats) GetId() string {
	return o.Id
}

func (o *ObjectFromNats) SetId(v string) {
	o.Id = v
}

func (o *ObjectFromNats) GetData() []byte {
	return o.Data
}

func (o *ObjectFromNats) SetData(v []byte) {
	o.Data = v
}

//--- для ObjectToNats ---

func (o *ObjectToNats) GetId() string {
	return o.Id
}

func (o *ObjectToNats) SetId(v string) {
	o.Id = v
}

func (o *ObjectToNats) GetData() any {
	return o.Data
}

func (o *ObjectToNats) SetData(v any) {
	o.Data = v
}

func (o *ObjectToNats) GetError() error {
	return o.Error
}

func (o *ObjectToNats) SetError(v error) {
	o.Error = v
}

func (o *ObjectToNats) GetTaskId() string {
	return o.TaskId
}

func (o *ObjectToNats) SetTaskId(v string) {
	o.TaskId = v
}

func (o *ObjectToNats) GetSource() string {
	return o.Source
}

func (o *ObjectToNats) SetSource(v string) {
	o.Source = v
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

// WithSubscription 'слушатель' запросов на поиск информации
func WithSubscription(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'subscription' cannot be empty")
		}

		n.subscriptionRequest = v

		return nil
	}
}
