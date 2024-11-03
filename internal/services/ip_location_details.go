package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/spf13/viper"
)

type Country struct {
	Cca2       string              `json:"cca2"`
	Currencies map[string]Currency `json:"currencies"`
	Languages  map[string]string   `json:"languages"`
	LatLng     []float64           `json:"latlng"`
	Name       CountryName         `json:"name"`
	Timezones  []string            `json:"timezones"`
}

type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CountryName struct {
	Common     string                 `json:"common"`
	Official   string                 `json:"official"`
	NativeName map[string]NativeNames `json:"nativeName"`
}

type NativeNames struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

type RatesData struct {
	Rates map[string]float64 `json:"rates"`
}

type IPLocationDetails struct {
	CountryName           string
	RelativeRates         map[string]float64
	CurrentTimeByTimezone map[string]string
	Currencies            map[string]Currency
	Cca2                  string
}

func GetIPLocationDetails(redisCache interfaces.Cache, httpClient interfaces.Client, ip string) (*IPLocationDetails, error) {
	ipLocation, err := NewIPLocation(httpClient)
	if err != nil {
		return nil, fmt.Errorf("error creating IP location service: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("context.timeout"))
	defer cancel()

	info, err := ipLocation.GetIPLocation(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("error getting IP location: %w", err)
	}

	country, err := getCountryFromCache(ctx, redisCache, info.IsoCode)
	if err != nil {
		return nil, err
	}

	ratesData, err := getRatesDataFromCache(ctx, redisCache)
	if err != nil {
		return nil, err
	}

	usdRate, exists := ratesData.Rates["USD"]
	if !exists {
		return nil, fmt.Errorf("USD rate not found in cache")
	}

	relativeRates := calculateRelativeRates(country.Currencies, ratesData.Rates, usdRate)

	currentTimeByTimezone := make(map[string]string)
	for _, timezone := range country.Timezones {
		offset, err := parseTimezoneOffset(timezone)
		if err == nil {
			currentTimeByTimezone[timezone] = time.Now().UTC().Add(offset).Format(time.RFC1123)
		}
	}

	return &IPLocationDetails{
		CountryName:           country.Name.Common,
		RelativeRates:         relativeRates,
		CurrentTimeByTimezone: currentTimeByTimezone,
		Currencies:            country.Currencies,
		Cca2:                  country.Cca2,
	}, nil

}

func getCountryFromCache(ctx context.Context, cache interfaces.Cache, isoCode string) (Country, error) {
	data, err := cache.Get(ctx, "countries")
	if err != nil || data == nil {
		return Country{}, fmt.Errorf("error getting or country data not found in cache")
	}

	var countries []Country
	if err := json.Unmarshal(data, &countries); err != nil {
		return Country{}, fmt.Errorf("error unmarshalling country data: %w", err)
	}

	for _, country := range countries {
		if country.Cca2 == isoCode {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("country code %s not found in cache", isoCode)
}

func getRatesDataFromCache(ctx context.Context, cache interfaces.Cache) (RatesData, error) {
	data, err := cache.Get(ctx, "currencies")
	if err != nil || data == nil {
		return RatesData{}, fmt.Errorf("error getting or currency data not found in cache")
	}

	var ratesData RatesData
	if err := json.Unmarshal(data, &ratesData); err != nil {
		return RatesData{}, fmt.Errorf("error unmarshalling currency data: %w", err)
	}

	return ratesData, nil
}

func calculateRelativeRates(currencies map[string]Currency, rates map[string]float64, usdRate float64) map[string]float64 {
	relativeRates := make(map[string]float64)
	for code := range currencies {
		if rate, exists := rates[code]; exists {
			relativeRates[code] = rate / usdRate
		}
	}
	return relativeRates
}

func parseTimezoneOffset(timezone string) (time.Duration, error) {
	if len(timezone) < 9 || timezone[:3] != "UTC" {
		return 0, fmt.Errorf("invalid timezone format")
	}

	sign := 1
	if timezone[3] == '-' {
		sign = -1
	}

	hours := 0
	minutes := 0
	_, err := fmt.Sscanf(timezone[4:], "%02d:%02d", &hours, &minutes)
	if err != nil {
		return 0, err
	}

	offset := time.Duration(sign) * (time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
	return offset, nil
}
