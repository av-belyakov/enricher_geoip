package geoipapi

import (
	"net/http"
	"time"
)

// GeoIpClient GeoIP клиента для запроса информации из БД GeoIP компании
type GeoIpClient struct {
	port              int
	host              string
	path              string
	connectionTimeout time.Duration
	client            *http.Client
}

type resultGeoIP struct {
	AddressVersion string          `json:"address_version"`
	IpLocations    []ipLocationSet `json:"ip_locations"`
}

type ipLocationSet struct {
	Source string `json:"source"`
	IpLocation
}

// GeoIpInformation список найденной информации по запрашиваемому ip адресу
type GeoIpInformation struct {
	IsSuccess bool
	Ip        string
	Info      map[string]IpLocation
}

// IpLocation подробная информация об ip адресе
type IpLocation struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

// geoIpClientOptions функциональные параметры
type geoIpClientOptions func(*GeoIpClient) error
