package dbip

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
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

func (d *Database) GetRequestsCountByIP(ctx context.Context, ip string) (int, error) {
	strVal, err := d.db.Get(ctx, ip).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, service.ErrNoSuchItem
		}

		return 0, fmt.Errorf("cannot execute get query: %w", err)
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, fmt.Errorf("cannot convert result into int: %w", err)
	}

	return val, nil
}
func (d *Database) SetAPIQuotaByIP(ctx context.Context, ip string, quota int, expiration time.Duration) error {
	err := d.db.Set(ctx, ip, quota, expiration).Err()
	if err != nil {
		return fmt.Errorf("cannot execute set query: %w", err)
	}

	return nil
}

func (d *Database) DecrAPIQuotaByIP(ctx context.Context, ip string) (int64, error) {
	remaining, err := d.db.Decr(ctx, ip).Result()
	if err != nil {
		return 0, fmt.Errorf("cannot execute decr query: %w", err)
	}

	return remaining, nil
}

func (d *Database) GetTTLByIP(ctx context.Context, ip string) (time.Duration, error) {
	ttl, err := d.db.TTL(ctx, ip).Result()
	if err != nil {
		return 0, fmt.Errorf("cannot execute ttl query: %w", err)
	}

	return ttl, nil
}

func (d *Database) IncrRequestCounter(ctx context.Context) error {
	err := d.db.Incr(ctx, "counter").Err()
	if err != nil {
		return fmt.Errorf("cannot execute incr query: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
