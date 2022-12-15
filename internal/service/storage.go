package service

import (
	"context"
	"time"
)

type Service struct {
	urlStorage URLStorage
	ipStorage  IPStorage
	apiQuota   int
}

type URLStorage interface {
	GetURLByID(ctx context.Context, id string) (string, error)
	InsertURL(ctx context.Context, id, urlStr string, ttl time.Duration) error
}
type IPStorage interface {
	IncrByIP(ctx context.Context, ip string) error
	CheckRateLimit(ctx context.Context, ip string) (string, error)
}
