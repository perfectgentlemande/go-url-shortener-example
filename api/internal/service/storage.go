package service

import "context"

type SlugStorage interface {
	GetByID(ctx context.Context, id string) (string, error)
}
