package service

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) Resolve(ctx context.Context, id, ip string) (string, error) {
	urlStr, err := s.urlStorage.GetURLByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get URL by ID: %w", err)
	}

	err = s.ipStorage.IncrByIP(ctx, ip)
	if err != nil {
		return "", fmt.Errorf("cannot incr by IP: %w", err)
	}

	return urlStr, nil
}

func (s *Service) GetURLByID(ctx context.Context, id string) (string, error) {

	return "", nil
}

func (s *Service) InsertURL(ctx context.Context, id, urlStr string, ttl time.Duration) error {
	return nil
}
func (s *Service) Incr(ctx context.Context)
