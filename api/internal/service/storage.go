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
	GetRequestsCountByIP(ctx context.Context, ip string) (string, error)
	SetRequestsCountByIP(ctx context.Context, id, url string, expiration time.Duration) error
}
