package dbip

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
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

func (d *Database) GetRequestsCountByIP(ctx context.Context, ip string) (string, error) {
	return "", nil
}
func (d *Database) SetRequestsCountByIP(ctx context.Context, id, url string, expiration time.Duration) error {
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
