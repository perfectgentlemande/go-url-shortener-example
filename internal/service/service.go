package service

import (
	"context"
	"time"
)

func (s *Service) GetURLByID(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (s *Service) InsertURL(ctx context.Context, id, urlStr string, ttl time.Duration) error {
	return nil
}
func (s *Service) Incr(ctx context.Context)
