package router

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/av-belyakov/enricher_geoip/cmd/geoipapi"
	"github.com/av-belyakov/enricher_geoip/cmd/natsapi"
	"github.com/av-belyakov/enricher_geoip/interfaces"
	"github.com/av-belyakov/enricher_geoip/internal/requests"
	"github.com/av-belyakov/enricher_geoip/internal/responses"
	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
)

type Router struct {
	counter       interfaces.Counter
	logger        interfaces.Logger
	geoIpClient   *geoipapi.GeoIpClient
	chFromNatsApi <-chan interfaces.Requester
	chToNatsApi   chan<- interfaces.Responser
}

func NewRouter(
	counter interfaces.Counter,
	logger interfaces.Logger,
	geoIpClient *geoipapi.GeoIpClient,
	chFromNatsApi <-chan interfaces.Requester,
	chToNatsApi chan<- interfaces.Responser,
) *Router {
	return &Router{
		geoIpClient:   geoIpClient,
		counter:       counter,
		logger:        logger,
		chFromNatsApi: chFromNatsApi,
		chToNatsApi:   chToNatsApi,
	}
}

func (r *Router) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-r.chFromNatsApi:
				go r.handlerRequest(ctx, msg)
			}
		}
	}()
}

func (r *Router) handlerRequest(ctx context.Context, msg interfaces.Requester) {
	response := &natsapi.ObjectToNats{
		Id: msg.GetId(),
	}

	if ctx.Err() != nil {
		return
	}

	var req requests.Request
	if err := json.Unmarshal(msg.GetData(), &req); err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		response.Error = errors.New("the request received an incorrect json format")
		r.chToNatsApi <- response

		return
	}

	response.TaskId = req.TaskId
	response.Source = req.Source

	results := make([]responses.DetailedInformation, 0, len(req.ListIp))
	for _, ip := range req.ListIp {
		result := responses.DetailedInformation{IpAddr: ip}

		res, err := r.geoIpClient.GetGeoInformation(ctx, ip)
		if err != nil {
			result.Error = "error interacting with a remote database"
			results = append(results, result)
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		var geoIPRes responses.ResponseGeoIPDataBase
		if err = json.Unmarshal(res, &geoIPRes); err != nil {
			result.Error = "a json object in an incorrect format was received from the geoip database"
			results = append(results, result)
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		geoIpInfo, _ := supportingfunctions.GetGeoIPInfo(geoIPRes)
		geoIpInfo.IpAddr = ip

		results = append(results, geoIpInfo)
	}

	response.Data = results

	r.counter.SendMessage("update processed events", 1)

	r.chToNatsApi <- response
}
