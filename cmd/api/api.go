package api

import (
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

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
