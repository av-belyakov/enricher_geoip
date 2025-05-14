package wrappers

// WrappersZabbixInteractionSettings настройки для обёртки взаимодействия с модулем zabbixapi
type WrappersZabbixInteractionSettings struct {
	EventTypes  []EventType //типы событий
	NetworkHost string      //ip адрес или доменное имя
	ZabbixHost  string      //zabbix host
	NetworkPort int         //сетевой порт
}

type Handshake struct {
	Message      string
	TimeInterval int
}

type EventType struct {
	EventType  string
	ZabbixKey  string
	Handshake  Handshake
	IsTransmit bool
}
