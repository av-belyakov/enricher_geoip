package geoipapi

import (
	"errors"
	"time"
)

// WithPort устанавливает порт для взаимодействия с модулем
func WithPort(v int) geoIpClientOptions {
	return func(gic *GeoIpClient) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		gic.port = v

		return nil
	}
}

// WithHost устанавливает хост для взаимодействия с модулем
func WithHost(v string) geoIpClientOptions {
	return func(gic *GeoIpClient) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		gic.host = v

		return nil
	}
}

// WithPath устанавливает путь запроса по которой осуществляется маршрутизация
func WithPath(v string) geoIpClientOptions {
	return func(gic *GeoIpClient) error {
		gic.path = v

		return nil
	}
}

// WithConnectionTimeout устанавливает время ожидания выполнения запроса
func WithConnectionTimeout(timeout time.Duration) geoIpClientOptions {
	return func(gic *GeoIpClient) error {
		if timeout > (1 * time.Second) {
			gic.connectionTimeout = timeout
		}

		return nil
	}
}
