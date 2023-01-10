package dburl

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"
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
	return Database{db: redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.No,
	})}, nil
}

func (d *Database) GetByID(ctx context.Context, id string) (string, error) {
	value, err := d.db.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return "", service.ErrNoSuchID
		}

		return "", fmt.Errorf("cannot get URL by ID: %w", err)
	}

	return value, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
