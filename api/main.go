package main

import (
	"log"

	"github.com/perfectgentlemande/go-url-shortener-example/api/database"
	"github.com/perfectgentlemande/go-url-shortener-example/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App, c *routes.Controller) {
	app.Get("/:url", c.Resolve)
	app.Post("/api/v1", c.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	r := database.CreateClient(0)
	defer r.Close()

	rInr := database.CreateClient(1)
	defer rInr.Close()

	// Implement Rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()

	c := &routes.Controller{
		R:    r,
		RInr: rInr,
		R2:   r2,
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app, c)
	app.Listen(":3000") // + os.Getenv("APP_PORT"))
	// base62EncodedString := helpers.Base62Encode(9999999)
	// fmt.Println(base62EncodedString)
	// fmt.Println(helpers.Base62Decode(base62EncodedString))
}
