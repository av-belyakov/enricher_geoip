package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/av-belyakov/enricher_geoip/constants"
	"github.com/av-belyakov/enricher_geoip/internal/appname"
	"github.com/av-belyakov/enricher_geoip/internal/appversion"
	"github.com/av-belyakov/enricher_geoip/internal/responses"
)

func getInformationMessage() string {
	version, err := appversion.GetVersion()
	if err != nil {
		log.Println(err)
	}

	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_ENRICHERGEOIP_MAIN")
	if ok && (envValue == "development") {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), strings.Replace(version, "\n", "", -1))

	fmt.Printf("\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", constants.Underlining, constants.Ansi_Bright_Green, appStatus, constants.Ansi_Reset)

	return msg
}

// GetInfoGeoIP возвращает информацию из списка найденных данны о геопозиционировании
// с самым высоким рейтингом, если рейтинг одинаковый то предпочтение отдается базе данных с
// заданным именем
func GetInfoGeoIP(data responses.ResponseGeoIPDataBase, preferredDB string) responses.DetailedInformation {
	/*
	   переделать поиск GeoIP
	   При укладывании в elastic необходимо определять geoip по следующему алгоритму:
	   1. если указан тег в кейсе - берем из тега ATs:geoip="Индия"
	   2. если нет: запрашиваем во внешней базе и берем самый высокий по рейтингу rating.
	   3. если рейтинг совпадает берем текущую приоритезацию.
	*/

	result := responses.DetailedInformation{}

	var rating int
	for _, info := range data.IpLocations {
		if rating < info.Rating {
			rating = info.Rating

			result.Code = info.CountryCode
			result.City = info.City
			result.Subnet = info.Subnet
			result.Country = info.Country
			result.UpdatedAt = info.UpdatedAt
			result.IpRange = struct {
				IpFirst string `json:"ip_first"`
				IpLast  string `json:"ip_last"`
			}{
				IpFirst: info.IpRange.IpFirst,
				IpLast:  info.IpRange.IpLast,
			}
		}
	}

	return result
}

func groupIpInfoResult(infoEvent datamodels.InformationFromEventEnricher) struct{ city, country, countryCode string } {
	sources := [...]string{"GeoipNoc", "MAXMIND", "DBIP", "AriadnaDB"}
	customIpResult := struct{ city, country, countryCode string }{}

	for _, ip := range infoEvent.GetIpAddresses() {
		for _, source := range sources {
			if city, ok := infoEvent.SearchCity(ip, source); ok && city != "" {
				if customIpResult.city != "" {
					continue
				}

				customIpResult.city = city
			}

			if country, ok := infoEvent.SearchCountry(ip, source); ok && country != "" {
				if customIpResult.country != "" {
					continue
				}

				customIpResult.country = country
			}

			if countryCode, ok := infoEvent.SearchCountryCode(ip, source); ok && countryCode != "" {
				if customIpResult.countryCode != "" {
					continue
				}

				customIpResult.countryCode = countryCode
			}
		}
	}

	return customIpResult
}
