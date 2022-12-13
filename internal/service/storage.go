package service

import (
	"context"
	"time"
)

type Service struct {
}

type URLStorage interface {
	GetURLByID(ctx context.Context, id string) (string, error)
	InsertURL(ctx context.Context, id, urlStr string, ttl time.Duration) error
}
type IncrStorage interface {
	Incr(ctx context.Context, counterKey string) error
}
