package dburl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
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

func (d *Database) GetURLByID(ctx context.Context, id string) (string, error) {
	urlStr, err := d.db.Get(ctx, id).Result()
	if err != nil {
		return "", fmt.Errorf("cannot get URL by id: %s: %w", id, err)
	}

	return urlStr, nil
}
func (d *Database) InsertURL(ctx context.Context, id, urlStr string, ttl time.Duration) error {
	err := d.db.Set(ctx, id, urlStr, ttl).Err()
	if err != nil {
		return fmt.Errorf("cannot insert URL: %s with id: %s: %w", urlStr, id, err)
	}

	return nil
}
