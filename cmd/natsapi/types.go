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
	chFromModule         chan ObjectForTransfer
	chToModule           chan ObjectForTransfer
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

// ObjectForTransfer объект для передачи данных
type ObjectForTransfer struct {
	Data   []byte
	Error  error
	TaskId string
}
