package service

import (
	"context"
	"time"
)

type URLStorage interface {
	GetByID(ctx context.Context, id string) (string, error)
	Set(ctx context.Context, id, url string, expiration time.Duration) error
}
