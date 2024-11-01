package api

import (
	"github.com/cgiraldoz/geo-ip-info/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func StartAPI(redisCache cache.Cache) {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		/*ctx := c.Context()
		redisService := cache.NewDefaultRedisService()
		client := cache.NewDefaultCache(redisService).NewClient()

		if err := client.Set(ctx, "key", "Cristian", 0).Err(); err != nil {
			panic(err)
		}

		val, err := client.Get(ctx, "key").Result()
		if err != nil {
			panic(err)
		}

		return c.SendString(val)
		//lite := geolite.NewGeoLite()
		//info := lite.GetLocation("103.103.184.1")
		//return c.JSON(info)*/
		println(redisCache)
		return c.SendString("Hello, World!")
	})

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
