package main

import (
	"context"
	"encoding/json"

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
	chFromNatsApi <-chan natsapi.ObjectForTransfer
	chToNatsApi   chan<- natsapi.ObjectForTransfer
}

func NewRouter(
	counter interfaces.Counter,
	logger interfaces.Logger,
	geoIpClient *geoipapi.GeoIpClient,
	chFrom <-chan natsapi.ObjectForTransfer,
	chTo chan<- natsapi.ObjectForTransfer,
) *Router {
	return &Router{
		counter:       counter,
		logger:        logger,
		chFromNatsApi: chFrom,
		chToNatsApi:   chTo,
	}
}

func (r *Router) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg := <-r.chFromNatsApi:
			go r.handlerRequest(ctx, msg)
		}
	}
}

func (r *Router) handlerRequest(ctx context.Context, msg natsapi.ObjectForTransfer) {
	if ctx.Err() != nil {
		return
	}

	var req requests.Request
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	results := make([]responses.DetailedInformation, 0, len(req.ListIp))
	for _, ip := range req.ListIp {
		result := responses.DetailedInformation{IpAddr: ip}

		res, err := r.geoIpClient.GetGeoInformation(ctx, ip)
		if err != nil {
			result.Error = err.Error()
			results = append(results, result)
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		var geoIPRes responses.ResponseGeoIPDataBase
		if err = json.Unmarshal(res, &geoIPRes); err != nil {
			result.Error = err.Error()
			results = append(results, result)
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		geoIpInfo := GetInfoGeoIP(geoIPRes)
	}
	r.counter.SendMessage("update processed events", 1)
}
