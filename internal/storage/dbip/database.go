package dburl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Database struct {
	db *redis.Client
}

func NewDatabase(conf *Config) *Database {
	return &Database{
		db: redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		}),
	}
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) Ping(ctx context.Context) error {
	return d.db.Ping(ctx).Err()
}

func (d *Database) IncrByIP(ctx context.Context, ip string) error {
	err := d.db.Incr(ctx, ip).Err()
	if err != nil {
		return fmt.Errorf("cannot incr ip: %s: %w", ip, err)
	}

	return nil
}

func (d *Database) CheckRateLimit(ctx context.Context, ip string, apiQuota int) (string, error) {
	val, err := d.db.Get(ctx, ip).Result()
	limit, _ := d.db.TTL(ctx, ip).Result()

	if err == redis.Nil {
		_ = d.db.Set(ctx, ip, strconv.Itoa(apiQuota), 30*60*time.Second).Err()
	} else if err == nil {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}
}
