package dbip

import (
	"context"
	"fmt"
	"strconv"
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

func NewDatabase(conf *Config) Database {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.No,
	})

	return Database{
		db: cli,
	}
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

func (d *Database) Ping(ctx context.Context) error {
	return d.db.Ping(ctx).Err()
}
func (d *Database) Close() error {
	return d.db.Close()
}

func ProvideStorage(conf *Config, lifecycle fx.Lifecycle) (service.IPStorage, error) {
	// Implement Rate limiting
	ipStorage := NewDatabase(conf)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ipStorage.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return ipStorage.Close()
		},
	})

	return &ipStorage, nil
}
