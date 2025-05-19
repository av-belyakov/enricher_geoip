package nats_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	NATS_HOST = "192.168.9.208"
	NATS_PORT = 4222

	CACHETTL     = 360
	SUBSCRIPTION = "object.geoip-request.test"
)

type ResponseData struct {
	Information []IPAddressInformation `json:"found_information"`
	Source      string                 `json:"source"`
	TaskId      string                 `json:"task_id"`
	Error       string                 `json:"error"`
}

type IPAddressInformation struct {
	IP      string `json:"ip_address"`
	Code    string `json:"code"`
	Country string `json:"country"`
	City    string `json:"city"`
	IpRange struct {
		IpFirst string `json:"ip_first"`
		IpLast  string `json:"ip_last"`
	} `json:"ip_range"`
	Subnet    string `json:"subnet"`
	UpdatedAt string `json:"updated_at"`
	Error     string `json:"error"`
}

func CreateNatsConnect(host string, port int) (*nats.Conn, error) {
	var (
		nc  *nats.Conn
		err error
	)

	nc, err = nats.Connect(
		fmt.Sprintf("%s:%d", host, port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return nc, err
	}

	fmt.Println("func 'CreateNatsConnect', START")

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		if err != nil {
			fmt.Println(err)
			fmt.Printf("func 'CreateNatsConnect' the connection with NATS has been disconnected %s\n", err.Error())

			return
		}

		fmt.Println("func 'CreateNatsConnect' the connection with NATS has been disconnected")
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		if err != nil {
			fmt.Printf("func 'CreateNatsConnect' the connection to NATS has been re-established %s\n", err.Error())

			return
		}

		fmt.Println("func 'CreateNatsConnect' the connection to NATS has been re-established")
	})

	return nc, nil
}

func TestGeoIpRequest(t *testing.T) {
	nc, err := CreateNatsConnect(NATS_HOST, NATS_PORT)
	if err != nil {
		log.Fatalln(err)
	}

	nmsg, err := nc.RequestWithContext(t.Context(), SUBSCRIPTION, []byte(`{
			"source": "test_source",
	  		"task_id": "dg87w82883r33r4qds",
	   		"list_ip_addresses": ["57.31.173.10", "71.67.123.36", "69.111.36.11"]
		}`))

	response := ResponseData{}
	err = json.Unmarshal(nmsg.Data, &response)
	assert.NoError(t, err)

	t.Log("Response:", response)

	assert.Equal(t, len(response.Information), 3)

	nmsg.Sub.Unsubscribe()

	t.Cleanup(func() {
		nc.Close()
	})
}
