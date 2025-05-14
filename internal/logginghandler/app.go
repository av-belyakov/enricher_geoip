package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/enricher_geoip/interfaces"
)

// New конструктор обработки логов
// это просто мост соединяющий несколько пакетов логирования,
// например, как в данном случае пакет simplyLogger и пакет взаимодействия с Zabbix
func New(writer interfaces.WriterLoggingData, chSysMonit chan<- interfaces.Messager) *LoggingChan {
	return &LoggingChan{
		dataWriter:           writer,
		chanSystemMonitoring: chSysMonit,
		chanLogging:          make(chan interfaces.Messager),
	}
}

// Start обработчик и распределитель логов
func (lc *LoggingChan) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-lc.chanLogging:
				//передача информации по логам в simpleLogger
				//здесь так же может быть вывод в консоль, с индикацией цветов соответствующих
				//определенному типу сообщений но для этого надо включить вывод на stdout
				//в конфигурационном файле
				_ = lc.dataWriter.Write(msg.GetType(), msg.GetMessage())

				//передача информации в некоторую систему мониторинга, например Zabbix
				if msg.GetType() == "error" || msg.GetType() == "warning" {
					msg := NewMessageLogging()
					msg.SetType("error")
					msg.SetMessage(fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()))

					lc.chanSystemMonitoring <- msg
				}
				if msg.GetType() == "info" {
					msg := NewMessageLogging()
					msg.SetType("info")
					msg.SetMessage(msg.GetMessage())

					lc.chanSystemMonitoring <- msg
				}
			}
		}
	}()
}
