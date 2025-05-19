package confighandler

// ConfigApp конфигурационные настройки приложения
type ConfigApp struct {
	Common  CfgCommon
	NATS    CfgNats
	LogDB   CfgWriteLogDB
	GeoIPDB CfgGeoIPDB
}

// CfgCommon общие настройки
type CfgCommon struct {
	Logs   []*LogSet
	Zabbix ZabbixOptions
}

// Logs настройки логирования
type Logs struct {
	Logging []*LogSet
}

type LogSet struct {
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
}

// ZabbixSet настройки доступа к Zabbix
type ZabbixSet struct {
	Zabbix ZabbixOptions
}

// ZabbixOptions настройки доступа к Zabbix
type ZabbixOptions struct {
	EventTypes  []EventType `yaml:"eventType"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
}

type EventType struct {
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
	IsTransmit bool      `yaml:"isTransmit"`
}

type Handshake struct {
	Message      string `validate:"required" yaml:"message"`
	TimeInterval int    `yaml:"timeInterval"`
}

// CfgNats настройки доступа к NATS
type CfgNats struct {
	Subscription string `yaml:"subscription"`
	Host         string `validate:"required" yaml:"host"`
	Port         int    `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL     int    `validate:"gt=10,lte=86400" yaml:"cache_ttl"`
}

// CfgWriteLogDB настройки записи данных в БД
type CfgWriteLogDB struct {
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	NameDB        string `yaml:"namedb"`
	StorageNameDB string `yaml:"storage_name_db"`
	Port          int    `validate:"gt=0,lte=65535" yaml:"port"`
}

// CfgGeoIPDB настройки взаимодействия с БД GeoIP
type CfgGeoIPDB struct {
	Host           string `yaml:"host"`
	Path           string `yaml:"path"`
	Port           int    `validate:"gt=0,lte=65535" yaml:"port"`
	RequestTimeout int    `validate:"gt=1,lt=13" yaml:"request_timeout"`
}
