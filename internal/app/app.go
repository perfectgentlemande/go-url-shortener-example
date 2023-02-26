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

func newServer() (*http.Server, error) {
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

	fmt.Println(conf)

	urlStorage, err := dburl.NewDatabase(context.TODO(), &dburl.Config{
		Addr:     conf.DBAddr,
		Password: conf.DBPass,
		No:       0,
	})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot create URL Storage: %w", err)
	}
	defer urlStorage.Close()

	// Implement Rate limiting
	ipStorage, err := dbip.NewDatabase(context.TODO(), &dbip.Config{
		Addr:     conf.DBAddr,
		Password: conf.DBPass,
		No:       1,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create IP Storage: %w", err)
	}
	defer ipStorage.Close()

	c := api.New(service.New(conf.APIQuota, &urlStorage, &ipStorage))
	r := chi.NewRouter()

	api.HandlerFromMux(c, r)

	return &http.Server{
		Handler: r,
		Addr:    conf.AppPort,
	}, nil
}

func registerHooks(lifecycle fx.Lifecycle, srv *http.Server) {
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
				return srv.Shutdown(ctx)
			},
		},
	)
}

var Module = fx.Options(
	fx.Provide(newServer),
	fx.Invoke(registerHooks),
)
