package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/api"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dbip"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dburl"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	DBAddr   string
	DBPass   string
	AppPort  string
	Domain   string
	APIQuota int
}

func provideConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load environment file: %w", err)
	}

	conf := Config{
		DBAddr:   viper.GetString("DB_ADDR"),
		DBPass:   viper.GetString("DB_PASS"),
		AppPort:  viper.GetString("APP_PORT"),
		APIQuota: viper.GetInt("API_QUOTA"),
	}

	return &conf, nil
}

func provideURLStorage(conf *Config) (*dburl.Database, error) {
	urlStorage, err := dburl.NewDatabase(context.TODO(), &dburl.Config{
		Addr:     conf.DBAddr,
		Password: conf.DBPass,
		No:       0,
	})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot create URL Storage: %w", err)
	}

	return &urlStorage, nil
}

func provideIPStorage(conf *Config) (*dbip.Database, error) {
	// Implement Rate limiting
	ipStorage, err := dbip.NewDatabase(context.TODO(), &dbip.Config{
		Addr:     conf.DBAddr,
		Password: conf.DBPass,
		No:       1,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create IP Storage: %w", err)
	}

	return &ipStorage, nil
}

func provideServer(conf *Config, dbURL *dburl.Database, dbIP *dbip.Database) (*http.Server, error) {
	c := api.New(service.New(conf.APIQuota, dbURL, dbIP))
	r := chi.NewRouter()

	api.HandlerFromMux(c, r)

	return &http.Server{
		Handler: r,
		Addr:    conf.AppPort,
	}, nil
}

func registerHooks(lifecycle fx.Lifecycle, srv *http.Server, dbURL *dburl.Database, dbIP *dbip.Database) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					log.Println("Listening on:", srv.Addr)
					err := srv.ListenAndServe()
					if err != nil {
						log.Println("Server listening error:", err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				defer dbURL.Close()
				defer dbIP.Close()

				return srv.Shutdown(ctx)
			},
		},
	)
}

var Module = fx.Options(
	fx.Provide(provideConfig),
	fx.Provide(provideURLStorage),
	fx.Provide(provideIPStorage),
	fx.Provide(provideServer),
	fx.Invoke(registerHooks),
)
