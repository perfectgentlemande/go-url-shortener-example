package dbip

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

func (d *Database) GetRequestsCountByIP(ctx context.Context, ip string) (int, error) {
	strVal, err := d.db.Get(ctx, ip).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, service.ErrNoSuchItem
		}

		return 0, fmt.Errorf("get query failed: %w", err)
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, fmt.Errorf("cannot convert result into int: %w", err)
	}

	return val, nil
}
func (d *Database) SetRequestsCountByIP(ctx context.Context, id, url string, expiration time.Duration) error {
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
