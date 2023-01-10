package service

import "context"

type URLStorage interface {
	GetByID(ctx context.Context, id string) (string, error)
}
