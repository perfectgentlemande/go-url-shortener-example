package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/database/dbip"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/database/dburl"
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

	addr, pass := os.Getenv("DB_ADDR"), os.Getenv("DB_PASS")
	urlStorage, err := dburl.NewDatabase(context.TODO(), &dburl.Config{
		Addr:     addr,
		Password: pass,
		No:       0,
	})
	if err != nil {
		log.Printf("cannot create URL Storage: %s\n", err)
		return
	}
	defer urlStorage.Close()

	// Implement Rate limiting
	ipStorage, err := dbip.NewDatabase(context.TODO(), &dbip.Config{
		Addr:     addr,
		Password: pass,
		No:       1,
	})
	if err != nil {
		log.Printf("cannot create IP Storage: %s\n", err)
		return
	}
	defer ipStorage.Close()

	defaultAPIQuotaStr := os.Getenv("API_QUOTA")
	defaultAPIQuota, err := strconv.Atoi(defaultAPIQuotaStr)
	if err != nil {
		log.Printf("wrong API_QUOTA: %d: %s", defaultAPIQuota, err)
		return
	}

	c := &routes.Controller{
		UrlStorage:      &urlStorage,
		IpStorage:       &ipStorage,
		DefaultAPIQuota: defaultAPIQuota,
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app, c)
	app.Listen(":3000") // + os.Getenv("APP_PORT"))
	// base62EncodedString := helpers.Base62Encode(9999999)
	// fmt.Println(base62EncodedString)
	// fmt.Println(helpers.Base62Decode(base62EncodedString))
}
