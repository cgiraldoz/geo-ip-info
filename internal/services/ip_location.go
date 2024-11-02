package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
	"io"
	"net"
	"net/http"
	"strings"
)

type CountryInfo struct {
	Name    string `maxminddb:"name"`
	IsoCode string `maxminddb:"iso_code"`
}

type IPLocation struct {
	httpClient interfaces.Client
}

type ExternalAPIResponse struct {
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
}

func NewIPLocation(httpClient interfaces.Client) (*IPLocation, error) {
	if viper.GetString("ipapi.url") == "" {
		return nil, errors.New("IPAPI URL not configured in Viper")
	}
	return &IPLocation{
		httpClient: httpClient,
	}, nil
}

func (g *IPLocation) GetIPLocation(ctx context.Context, ip string) (*CountryInfo, error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb")

	if err != nil {
		return nil, fmt.Errorf("error opening local database: %w", err)
	}

	defer func(db *geoip2.Reader) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing local database")
		}
	}(db)

	record, err := db.City(net.ParseIP(ip))

	if err == nil && record != nil && record.Country.IsoCode != "" {
		return &CountryInfo{
			Name:    record.Country.Names["en"],
			IsoCode: record.Country.IsoCode,
		}, nil
	}

	return g.fetchFromAPI(ctx, ip)
}

func (g *IPLocation) fetchFromAPI(ctx context.Context, ip string) (*CountryInfo, error) {
	apiURL := viper.GetString("ipapi.url")
	if apiURL == "" {
		return nil, errors.New("IPAPI URL not configured")
	}

	url := strings.Replace(apiURL, "{ip}", ip, 1)
	resp, err := g.httpClient.Get(ctx, url)

	if err != nil {
		return nil, fmt.Errorf("error fetching IP location from external API: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	var apiResponse ExternalAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding API response: %w", err)
	}

	if apiResponse.CountryName == "" || apiResponse.CountryCode == "" {
		return nil, fmt.Errorf("IP location not found for IP: %s", ip)
	}

	return &CountryInfo{
		Name:    apiResponse.CountryName,
		IsoCode: apiResponse.CountryCode,
	}, nil
}
