package service

import (
	"context"
	"time"
)

type URLStorage interface {
	GetByID(ctx context.Context, id string) (string, error)
	SetByID(ctx context.Context, id, url string, expiration time.Duration) error
}

type IPStorage interface {
	GetRequestsCountByIP(ctx context.Context, ip string) (int, error)
	SetAPIQuotaByIP(ctx context.Context, ip string, quota int, expiration time.Duration) error
	DecrAPIQuotaByIP(ctx context.Context, ip string) (int64, error)
	GetTTLByIP(ctx context.Context, ip string) (time.Duration, error)
	IncrRequestCounter(ctx context.Context) error
}
