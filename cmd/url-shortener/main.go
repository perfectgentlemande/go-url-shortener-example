package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/api"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dbip"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dburl"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"

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

	defaultAPIQuota := viper.GetInt("API_QUOTA")
	c := api.New(service.New(defaultAPIQuota, &urlStorage, &ipStorage))
	r := chi.NewRouter()

	api.HandlerFromMux(c, r)

	s := &http.Server{
		Handler: r,
		Addr:    viper.GetString("APP_PORT"),
	}

	log.Fatal(s.ListenAndServe())
}
