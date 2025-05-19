package geoip_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/av-belyakov/enricher_geoip/cmd/geoipapi"
	"github.com/av-belyakov/enricher_geoip/cmd/natsapi"
	"github.com/av-belyakov/enricher_geoip/cmd/router"
	"github.com/av-belyakov/enricher_geoip/interfaces"
	"github.com/av-belyakov/enricher_geoip/internal/logginghandler"
	"github.com/av-belyakov/enricher_geoip/internal/responses"
	"github.com/stretchr/testify/assert"
)

type Counting struct {
	ch chan struct {
		msg string
		num int
	}
}

func (c *Counting) SendMessage(value string, count int) {
	c.ch <- struct {
		msg string
		num int
	}{
		msg: value,
		num: count,
	}
}

type Logging struct {
	ch chan interfaces.Messager
}

func (l *Logging) GetChan() <-chan interfaces.Messager {
	return l.ch
}

func (l *Logging) Send(msgType, msgData string) {
	l.ch <- &logginghandler.MessageLogging{
		Message: msgData,
		Type:    msgType,
	}
}

func TestRequestGeoIP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	geoIpClient, err := geoipapi.NewGeoIpClient(
		geoipapi.WithHost("pg2.cloud.gcm"),
		geoipapi.WithPort(88),
		geoipapi.WithPath("ip"),
		geoipapi.WithConnectionTimeout(time.Duration(7*time.Second)))
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("geoIpClient: '%+v'\n", geoIpClient)

	counting := &Counting{ch: make(chan struct {
		msg string
		num int
	})}
	logging := &Logging{ch: make(chan interfaces.Messager)}

	//счётчик
	go func(c *Counting) {
		for message := range c.ch {
			t.Log("counting: ", message)
		}
	}(counting)

	//логирование
	go func(l *Logging) {
		for message := range l.ch {
			t.Log("logging: ", message)
		}
	}(logging)

	chToRoute := make(chan interfaces.Requester)
	chFromRoute := make(chan interfaces.Responser)

	r := router.NewRouter(counting, logging, geoIpClient, chToRoute, chFromRoute)
	r.Start(ctx)

	taskId := time.Now().Format("20250515-125612.00000")

	chToRoute <- &natsapi.ObjectFromNats{
		Id: taskId,
		Data: fmt.Appendf(nil, `{
			"source": "test_script",
			"task_id": "%s",
			"list_ip_addresses": ["45.123.3.66", "78.99.100.3", "32.0.26.33"]
		}`, taskId),
	}

	data := <-chFromRoute

	res, ok := data.GetData().([]responses.DetailedInformation)
	assert.True(t, ok)
	assert.Equal(t, len(res), 3)

	t.Logf("%+v", res)

	t.Cleanup(func() {
		cancel()

		close(counting.ch)
		close(logging.ch)
	})
}
