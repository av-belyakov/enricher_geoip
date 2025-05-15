package natsapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
)

// subscriptionHandler обработчик подписки приёма запросов
func (api *apiNatsModule) subscriptionRequestHandler() {
	_, err := api.natsConn.Subscribe(api.subscriptionRequest, func(m *nats.Msg) {
		taskId := uuid.NewString()

		api.storage.SetReq(taskId, m)

		api.chFromModule <- ObjectForTransfer{
			TaskId: taskId,
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
			if m, ok := api.storage.GetReq(incomingData.TaskId); ok {
				m.Respond(incomingData.Data)
				api.storage.DelReq(incomingData.TaskId)

				continue
			}
		}
	}
}
