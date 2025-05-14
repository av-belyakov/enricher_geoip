package main

import (
	"context"
	"log"
	"os"

	"github.com/av-belyakov/enricher_geoip/cmd/elasticsearchapi"
	"github.com/av-belyakov/enricher_geoip/cmd/natsapi"
	"github.com/av-belyakov/enricher_geoip/cmd/wrappers"
	"github.com/av-belyakov/enricher_geoip/constants"
	"github.com/av-belyakov/enricher_geoip/interfaces"
	"github.com/av-belyakov/enricher_geoip/internal/confighandler"
	"github.com/av-belyakov/enricher_geoip/internal/countermessage"
	"github.com/av-belyakov/enricher_geoip/internal/logginghandler"
	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
	"github.com/av-belyakov/simplelogger"
)

func app(ctx context.Context) {
	var nameRegionalObject string
	if os.Getenv("GO_ENRICHERGEOIP_MAIN") == "development" {
		nameRegionalObject = "enricher_geoip-dev"
	} else {
		nameRegionalObject = "enricher_geoip"
	}

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// ****************************************************************************
	// *********** инициализируем модуль чтения конфигурационного файла ***********
	conf, err := confighandler.New(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	var listLog []simplelogger.OptionsManager
	for _, v := range conf.GetListLogs() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	confDB := conf.GetLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: nameRegionalObject,
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	// ************************************************************************
	// ************* инициализация модуля взаимодействия с Zabbix *************
	chZabbix := make(chan interfaces.Messager)
	confZabbix := conf.GetZabbix()
	wziSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confZabbix.NetworkPort,
		NetworkHost: confZabbix.NetworkHost,
		ZabbixHost:  confZabbix.ZabbixHost,
	}
	eventTypes := []wrappers.EventType(nil)
	for _, v := range confZabbix.EventTypes {
		eventTypes = append(eventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	wziSettings.EventTypes = eventTypes
	wrappers.WrappersZabbixInteraction(ctx, wziSettings, simpleLogger, chZabbix)

	//***************************************************************************
	//************** инициализация обработчика логирования данных ***************
	//фактически это мост между simpleLogger и пакетом соединения с Zabbix
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	// ***************************************************************************
	// *********** инициализируем модуль счётчика для подсчёта сообщений *********
	counting := countermessage.New(chZabbix)
	counting.Start(ctx)

	// ***********************************************************************
	// ************** инициализация модуля взаимодействия с NATS *************
	confNats := conf.NATS
	apiNats, err := natsapi.New(
		counting,
		logging,
		natsapi.WithHost(confNats.Host),
		natsapi.WithPort(confNats.Port),
		natsapi.WithCacheTTL(confNats.CacheTTL),
		natsapi.WithSubscription(confNats.Subscription))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}
	//--- старт модуля
	if err = apiNats.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

}
