package dburl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"go.uber.org/fx"
)

type Config struct {
	Addr     string
	Password string
	No       int
}

type Database struct {
	db *redis.Client
}

func NewDatabase(ctx context.Context, conf *Config) (Database, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.No,
	})

	err := cli.Ping(ctx).Err()
	if err != nil {
		return Database{}, fmt.Errorf("cannot ping database: %w", err)
	}

	return Database{
		db: cli,
	}, nil
}

func (d *Database) GetByID(ctx context.Context, id string) (string, error) {
	value, err := d.db.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return "", service.ErrNoSuchItem
		}

		return "", fmt.Errorf("get query failed: %w", err)
	}

	return value, nil
}

func (d *Database) SetByID(ctx context.Context, id, url string, expiration time.Duration) error {
	err := d.db.Set(ctx, id, url, expiration).Err()
	if err != nil {
		return fmt.Errorf("set query failed: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func ProvideStorage(ctx context.Context, conf *Config, lifecycle fx.Lifecycle) (service.URLStorage, error) {
	urlStorage, err := NewDatabase(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("cannot create URL Storage: %w", err)
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("url storage closed")
			return urlStorage.Close()
		},
	})

	return &urlStorage, nil
}
