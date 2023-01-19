package service

import (
	"context"
	"fmt"
)

type Service struct {
	defaultAPIQuota int
	urlStorage      URLStorage
	ipStorage       IPStorage
}

func New(defaultAPIQuota int, urlStorage URLStorage, ipStorage IPStorage) *Service {
	return &Service{
		defaultAPIQuota: defaultAPIQuota,
		urlStorage:      urlStorage,
		ipStorage:       ipStorage,
	}
}

func (s *Service) Resolve(ctx context.Context, id string) (string, error) {
	value, err := s.urlStorage.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get url by id: %w", err)
	}

	err = s.ipStorage.IncrRequestCounter(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot increment request counter: %w", err)
	}

	return value, nil
}
