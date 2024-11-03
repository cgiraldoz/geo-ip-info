package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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
	LatLng                []float64
	DistanceToBuenosAires float64
}

type DistanceStats struct {
	FarthestDistance float64
	ClosestDistance  float64
	TotalDistance    float64
	TotalRequests    int
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

	countryCacheKey := "country:" + info.IsoCode
	cachedDetails, err := getCountryDetailsFromCache(ctx, redisCache, countryCacheKey)
	if err == nil && cachedDetails != nil {
		updateCurrentTimeByTimezone(cachedDetails)
		cachedDetails.DistanceToBuenosAires = calculateDistanceToBuenosAires(cachedDetails.LatLng)
		updateDistanceStats(ctx, redisCache, cachedDetails.LatLng)
		return cachedDetails, nil
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
	currentTimeByTimezone := calculateCurrentTimeByTimezone(country.Timezones)

	ipDetails := &IPLocationDetails{
		CountryName:           country.Name.Common,
		RelativeRates:         relativeRates,
		CurrentTimeByTimezone: currentTimeByTimezone,
		Currencies:            country.Currencies,
		Cca2:                  country.Cca2,
		LatLng:                country.LatLng,
		DistanceToBuenosAires: calculateDistanceToBuenosAires(country.LatLng),
	}

	err = cacheCountryDetails(ctx, redisCache, countryCacheKey, ipDetails)
	if err != nil {
		return nil, fmt.Errorf("error caching country details: %w", err)
	}
	updateDistanceStats(ctx, redisCache, ipDetails.LatLng)
	return ipDetails, nil
}

func getCountryDetailsFromCache(ctx context.Context, cache interfaces.Cache, countryCacheKey string) (*IPLocationDetails, error) {
	data, err := cache.Get(ctx, countryCacheKey)
	if err != nil || data == nil {
		return nil, err
	}

	var details IPLocationDetails
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, fmt.Errorf("error unmarshalling country details from cache: %w", err)
	}

	return &details, nil
}

func cacheCountryDetails(ctx context.Context, cache interfaces.Cache, countryCacheKey string, details *IPLocationDetails) error {
	data, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("error marshalling country details for cache: %w", err)
	}

	ttl := viper.GetDuration("cache.ip_location_details.ttl")
	return cache.Set(ctx, countryCacheKey, data, ttl)
}

func updateCurrentTimeByTimezone(details *IPLocationDetails) {
	for timezone := range details.CurrentTimeByTimezone {
		offset, err := parseTimezoneOffset(timezone)
		if err == nil {
			details.CurrentTimeByTimezone[timezone] = time.Now().UTC().Add(offset).Format(time.RFC1123)
		}
	}
}

func calculateCurrentTimeByTimezone(timezones []string) map[string]string {
	currentTimeByTimezone := make(map[string]string)
	for _, timezone := range timezones {
		offset, err := parseTimezoneOffset(timezone)
		if err == nil {
			currentTimeByTimezone[timezone] = time.Now().UTC().Add(offset).Format(time.RFC1123)
		}
	}
	return currentTimeByTimezone
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

func updateDistanceStats(ctx context.Context, cache interfaces.Cache, countryLatLng []float64) {
	if len(countryLatLng) < 2 {
		fmt.Println("Error: countryLatLng does not contain valid coordinates")
		return
	}

	buenosAiresLat, buenosAiresLng, err := getBuenosAiresLatLng()
	if err != nil {
		fmt.Printf("Error getting Buenos Aires coordinates: %v\n", err)
		return
	}

	distance := calculateDistance(buenosAiresLat, buenosAiresLng, countryLatLng[0], countryLatLng[1])

	stats, err := getDistanceStatsFromCache(ctx, cache)
	if err != nil {
		stats = &DistanceStats{
			FarthestDistance: distance,
			ClosestDistance:  distance,
			TotalDistance:    distance,
			TotalRequests:    1,
		}
	} else {
		stats.TotalRequests++
		stats.TotalDistance += distance
		if distance > stats.FarthestDistance {
			stats.FarthestDistance = distance
		}
		if distance < stats.ClosestDistance || stats.ClosestDistance == 0 {
			stats.ClosestDistance = distance
		}
	}

	err = cacheDistanceStats(ctx, cache, stats)
	if err != nil {
		fmt.Printf("Error caching distance stats: %v\n", err)
	}
}

func getDistanceStatsFromCache(ctx context.Context, cache interfaces.Cache) (*DistanceStats, error) {
	data, err := cache.Get(ctx, "distance_stats")
	if err != nil || data == nil {
		return nil, err
	}

	var stats DistanceStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("error unmarshalling distance stats from cache: %w", err)
	}

	return &stats, nil
}

func cacheDistanceStats(ctx context.Context, cache interfaces.Cache, stats *DistanceStats) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("error marshalling distance stats for cache: %w", err)
	}

	return cache.Set(ctx, "distance_stats", data, viper.GetDuration("cache.distance_stats.ttl"))
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
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

func getBuenosAiresLatLng() (float64, float64, error) {
	lat := viper.GetFloat64("fixed_location.argentina.latitude")
	lng := viper.GetFloat64("fixed_location.argentina.longitude")

	if lat == 0 && lng == 0 {
		return 0, 0, fmt.Errorf("buenos Aires lat/lng not found or invalid in configuration")
	}

	return lat, lng, nil
}

func calculateDistanceToBuenosAires(countryLatLng []float64) float64 {
	if len(countryLatLng) < 2 {
		return 0
	}

	buenosAiresLat, buenosAiresLng, err := getBuenosAiresLatLng()
	if err != nil {
		fmt.Printf("error getting Buenos Aires coordinates: %v\n", err)
		return 0
	}

	return calculateDistance(buenosAiresLat, buenosAiresLng, countryLatLng[0], countryLatLng[1])
}
