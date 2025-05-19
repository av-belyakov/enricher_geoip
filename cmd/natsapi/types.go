package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_geoip/cmd/natsapi/storagetemporary"
	"github.com/av-belyakov/enricher_geoip/interfaces"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	counter              interfaces.Counter
	logger               interfaces.Logger
	natsConn             *nats.Conn
	storage              *storagetemporary.StorageTemporary
	subscriptionRequest  string
	subscriptionResponse string
	settings             apiNatsSettings
	chFromModule         chan interfaces.Requester
	chToModule           chan interfaces.Responser
}

type apiNatsSettings struct {
	nameRegionalObject string
	command            string
	host               string
	cachettl           int
	port               int
}

// NatsApiOptions функциональные опции
type NatsApiOptions func(*apiNatsModule) error

// ObjectFromNats объект для передачи данных
type ObjectFromNats struct {
	Data []byte
	Id   string
}

// ObjectToNats объект для передачи данных
type ObjectToNats struct {
	Data   any
	Error  error
	Id     string
	TaskId string
	Source string
}
