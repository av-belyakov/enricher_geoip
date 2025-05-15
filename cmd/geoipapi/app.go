package geoipapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
)

// NewGeoIpClient GeoIP клиент
func NewGeoIpClient(opts ...geoIpClientOptions) (*GeoIpClient, error) {
	settings := GeoIpClient{connectionTimeout: 30 * time.Second}

	for _, opt := range opts {
		if err := opt(&settings); err != nil {
			return &settings, err
		}
	}

	settings.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     settings.connectionTimeout,
			MaxIdleConnsPerHost: 10,
		}}

	return &settings, nil
}

// GetGeoInformation делает запрос к API БД GeoIP
func (gic *GeoIpClient) GetGeoInformation(ctx context.Context, ip string) ([]byte, error) {
	rex := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)[.]){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	tmp := rex.FindStringSubmatch(ip)
	if len(tmp) == 0 {
		return []byte{}, supportingfunctions.CustomError(fmt.Errorf("an invalid ip address '%s' was received", ip))
	}

	url := fmt.Sprintf("http://%s:%d/%s/%s/", gic.host, gic.port, gic.path, tmp[0])
	req, err := http.NewRequestWithContext(ctx, "GET", url, strings.NewReader(""))
	if err != nil {
		return []byte{}, supportingfunctions.CustomError(err)
	}

	res, err := gic.client.Do(req)
	defer supportingfunctions.ResponseClose(res)
	if err != nil {
		return []byte{}, supportingfunctions.CustomError(err)
	}

	if res.StatusCode != http.StatusOK {
		return []byte{}, supportingfunctions.CustomError(fmt.Errorf("error sending the request, response status is %s", res.Status))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resBody, supportingfunctions.CustomError(err)
	}

	return resBody, nil
}
