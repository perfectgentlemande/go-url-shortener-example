package main

import (
	"context"
	"log"
	"strconv"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/api"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dbip"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dburl"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	addr, pass := viper.GetString("DB_ADDR"), viper.GetString("DB_PASS")
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

	defaultAPIQuotaStr := viper.GetString("API_QUOTA")
	defaultAPIQuota, err := strconv.Atoi(defaultAPIQuotaStr)
	if err != nil {
		log.Printf("wrong API_QUOTA: %d: %s", defaultAPIQuota, err)
		return
	}

	c := api.New(service.New(defaultAPIQuota, &urlStorage, &ipStorage))
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/:url", c.Resolve)
	app.Post("/api/v1", c.Shorten)
	app.Listen(":3000") // + os.Getenv("APP_PORT"))
	// base62EncodedString := helpers.Base62Encode(9999999)
	// fmt.Println(base62EncodedString)
	// fmt.Println(helpers.Base62Decode(base62EncodedString))
}
