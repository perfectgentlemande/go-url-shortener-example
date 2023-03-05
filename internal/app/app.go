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
	DBURL   *dburl.Config
	DBIP    *dbip.Config
	API     *api.Config
	Service *service.Config
}

func provideAllConfigs() (*dburl.Config, *dbip.Config, *api.Config, *service.Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("could not load environment file: %w", err)
	}

	conf := Config{
		DBURL: &dburl.Config{
			Addr:     viper.GetString("DB_ADDR"),
			Password: viper.GetString("DB_PASS"),
			No:       0,
		},
		DBIP: &dbip.Config{
			Addr:     viper.GetString("DB_ADDR"),
			Password: viper.GetString("DB_PASS"),
			No:       1,
		},
		API: &api.Config{
			AppPort: viper.GetString("APP_PORT"),
			Domain:  viper.GetString("DOMAIN"),
		},
		Service: &service.Config{
			APIQuota: viper.GetInt("API_QUOTA"),
		},
	}

	return conf.DBURL, conf.DBIP, conf.API, conf.Service, nil
}

func provideServer(apiConf *api.Config, srvcConf *service.Config, dbURL *dburl.Database, dbIP *dbip.Database) (*http.Server, error) {
	c := api.New(service.New(srvcConf.APIQuota, dbURL, dbIP), apiConf.Domain)
	r := chi.NewRouter()

	api.HandlerFromMux(c, r)

	return &http.Server{
		Handler: r,
		Addr:    apiConf.AppPort,
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
				return srv.Shutdown(ctx)
			},
		},
	)
}

var Module = fx.Options(
	fx.Provide(provideAllConfigs),
	fx.Provide(dburl.ProvideStorage),
	fx.Provide(dbip.ProvideStorage),
	fx.Provide(provideServer),
	fx.Invoke(registerHooks),
)
