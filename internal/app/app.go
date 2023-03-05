package app

import (
	"fmt"

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

var Module = fx.Options(
	fx.Provide(provideAllConfigs),
	fx.Provide(dburl.ProvideStorage),
	fx.Provide(dbip.ProvideStorage),
	fx.Provide(service.Provide),
	fx.Invoke(api.Provide),
)
