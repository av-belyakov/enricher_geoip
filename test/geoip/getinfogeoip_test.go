package geoip_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/av-belyakov/enricher_geoip/internal/responses"
	"github.com/av-belyakov/enricher_geoip/internal/supportingfunctions"
	"github.com/stretchr/testify/assert"
)

func TestGetInfoGeoIp(t *testing.T) {
	f, err := os.ReadFile("../test_json/geoip_info.json")
	if err != nil {
		log.Fatal(err)
	}

	var geoIPRes responses.ResponseGeoIPDataBase
	err = json.Unmarshal(f, &geoIPRes)
	assert.NoError(t, err)

	result, source := supportingfunctions.GetGeoIPInfo(geoIPRes)
	assert.Equal(t, source, "GeoipNoc")
	assert.Equal(t, result.Code, "US")
	assert.Equal(t, result.Country, "США")
	assert.Equal(t, result.City, "United States")
}
