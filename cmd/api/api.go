package api

import (
	"github.com/cgiraldoz/geo-ip-info/internal/geolite"
	"github.com/gofiber/fiber/v2"
)

func StartAPI() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		lite := geolite.GeoLite{}
		info := lite.GetLocation("186.154.207.219")
		return c.JSON(info)
	})

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
