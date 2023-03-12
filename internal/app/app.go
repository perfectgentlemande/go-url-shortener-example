package app

import (
	"fmt"
	"os"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/api"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dbip"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/database/dburl"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/logger"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	DBURL   *dburl.Config
	DBIP    *dbip.Config
	API     *api.Config
	Service *service.Config
	Logger  *logger.Config
}

func provideAllConfigs() (*dburl.Config, *dbip.Config, *api.Config, *service.Config, *logger.Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("could not load environment file: %w", err)
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
		Logger: &logger.Config{
			Encoder: viper.GetString("LOG_ENCODER"),
			Level:   viper.GetString("LOG_LEVEL"),
		},
	}

	return conf.DBURL, conf.DBIP, conf.API, conf.Service, conf.Logger, nil
}

// I don't know how to insert anything apart from zap into this part
func addDefaultJSONLoggerForGoFx() fxevent.Logger {
	return &fxevent.ZapLogger{Logger: zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), os.Stdout, zap.DebugLevel))}
}

var Module = fx.Options(
	fx.WithLogger(addDefaultJSONLoggerForGoFx),
	fx.Provide(provideAllConfigs),
	fx.Provide(logger.ProvideLogger),
	fx.Provide(dburl.ProvideStorage),
	fx.Provide(dbip.ProvideStorage),
	fx.Provide(service.Provide),
	fx.Invoke(api.Provide),
)
