package api

import (
	"context"
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/gofiber/fiber/v2"
)

type IPDetails struct {
	CountryName           string                       `json:"country_name"`
	Cca2                  string                       `json:"cca2"`
	Currencies            map[string]services.Currency `json:"currencies"`
	RelativeRates         map[string]float64           `json:"relative_rates"`
	CurrentTimeByTimezone map[string]string            `json:"current_time_by_timezone"`
	DistanceToBuenosAires float64                      `json:"distance_to_buenos_aires"`
}

type DistanceStatsResponse struct {
	FarthestDistance float64                        `json:"farthest_distance"`
	FarthestCountry  string                         `json:"farthest_country"`
	ClosestDistance  float64                        `json:"closest_distance"`
	ClosestCountry   string                         `json:"closest_country"`
	TotalDistance    float64                        `json:"total_distance"`
	TotalRequests    int                            `json:"total_requests"`
	AverageDistance  float64                        `json:"average_distance"`
	CountryDistances map[string]CountryDistanceData `json:"country_distances"`
}

type CountryDistanceData struct {
	TotalDistance float64 `json:"total_distance"`
	Requests      int     `json:"requests"`
}

func StartAPI(redisCache interfaces.Cache, httpClient interfaces.Client) {
	app := fiber.New()

	app.Get("/api/ip/:ip", func(c *fiber.Ctx) error {
		ipDetails, err := services.GetIPLocationDetails(redisCache, httpClient, c.Params("ip"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(IPDetails{
			CountryName:           ipDetails.CountryName,
			Cca2:                  ipDetails.Cca2,
			Currencies:            ipDetails.Currencies,
			RelativeRates:         ipDetails.RelativeRates,
			CurrentTimeByTimezone: ipDetails.CurrentTimeByTimezone,
			DistanceToBuenosAires: ipDetails.DistanceToBuenosAires,
		})
	})

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		stats, err := services.GetDistanceStatsFromCache(context.Background(), redisCache)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving stats",
			})
		}

		averageDistance := services.CalculateWeightedAverageDistance(stats)

		countryDistances := make(map[string]CountryDistanceData)
		for country, data := range stats.CountryDistances {
			countryDistances[country] = CountryDistanceData{
				TotalDistance: data.TotalDistance,
				Requests:      data.Requests,
			}
		}

		return c.JSON(DistanceStatsResponse{
			FarthestDistance: stats.FarthestDistance,
			FarthestCountry:  stats.FarthestCountryName,
			ClosestDistance:  stats.ClosestDistance,
			ClosestCountry:   stats.ClosestCountryName,
			TotalDistance:    stats.TotalDistance,
			TotalRequests:    stats.TotalRequests,
			AverageDistance:  averageDistance,
			CountryDistances: countryDistances,
		})
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
