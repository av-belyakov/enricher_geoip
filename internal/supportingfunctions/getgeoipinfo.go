package supportingfunctions

import (
	"strconv"

	"github.com/av-belyakov/enricher_geoip/internal/responses"
)

// GetInfoGeoIP возвращает информацию из списка найденных данны о геопозиционировании
// Выбираются данные по следующему принципу:
//  1. приоритет отдаётся объекту с самым высоким рейтингом
//  2. если рейтинг у двух объектов одинаков то в приоритете источник
//     данных, сначала 'GeoipNoc', затем 'MAXMIND'
func GetGeoIPInfo(data responses.ResponseGeoIPDataBase) (responses.DetailedInformation, string, error) {
	var (
		result responses.DetailedInformation = responses.DetailedInformation{}
		rating int
		source string
	)

	for _, info := range data.IpLocations {
		r, err := strconv.Atoi(info.Rating)
		if err != nil {
			return result, source, err
		}

		if info.Country == "" || info.CountryCode == "" {
			continue
		}

		if rating > r {
			continue
		} else if rating == r {
			if source == "GeoipNoc" {
				continue
			} else if source == "MAXMIND" && info.Source != "GeoipNoc" {
				continue
			} else {
				rating = r
				source = info.Source

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
		} else {
			rating = r
			source = info.Source

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

	return result, source, nil
}
