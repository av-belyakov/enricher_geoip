package geoipapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
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
func (gic *GeoIpClient) GetGeoInformation(ctx context.Context, ip string) (GeoIpInformation, error) {
	result := GeoIpInformation{
		Ip:   ip,
		Info: make(map[string]IpLocation, 0),
	}

	rex := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)[.]){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	tmp := rex.FindStringSubmatch(ip)
	if len(tmp) == 0 {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("an invalid ip address '%s' was received %s:%d", ip, f, l-1)
	}

	url := fmt.Sprintf("http://%s:%d/%s/%s/", gic.host, gic.port, gic.path, tmp[0])
	req, err := http.NewRequestWithContext(ctx, "GET", url, strings.NewReader(""))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("%v %s:%d", err, f, l-2)
	}

	res, err := gic.client.Do(req)
	defer supportingfunctions.ResponseClose(res)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("%v %s:%d", err, f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("error sending the request, response status is %s %s:%d", res.Status, f, l-1)
	}

	resultGeoIP := resultGeoIP{}
	err = json.NewDecoder(res.Body).Decode(&resultGeoIP)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("%v %s:%d", err, f, l-2)
	}

	result.IsSuccess = true
	for _, v := range resultGeoIP.IpLocations {
		result.Info[v.Source] = IpLocation{
			City:        v.City,
			Country:     v.Country,
			CountryCode: v.CountryCode,
		}
	}

	return result, nil
}
