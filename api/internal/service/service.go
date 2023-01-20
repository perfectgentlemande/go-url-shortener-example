package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/api/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/base62"
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

func (s *Service) Shorten(ctx context.Context, ip, url, customShort string, expiry time.Duration) (string, int64, error) {
	limit, err := s.ipStorage.GetTTLByIP(ctx, ip)
	if err != nil {
		return "", 0, fmt.Errorf("cannot get TTL by IP: %w", err)
	}

	valInt, err := s.ipStorage.GetRequestsCountByIP(ctx, ip)
	if errors.Is(err, ErrNoSuchItem) {
		err = s.ipStorage.SetAPIQuotaByIP(ctx, ip, s.defaultAPIQuota, 30*60*time.Second)
		if err != nil {
			return "", 0, fmt.Errorf("cannot set api quota for IP: %w", err)
		}
	} else if err == nil {
		if valInt <= 0 {
			return "", 0, ErrRateLimitExceeded
		}
	}

	// enforce HTTPS, SSL
	url = helpers.EnforceHTTP(url)

	var id string
	if customShort == "" {
		id = base62.Encode(rand.Uint64())
	} else {
		id = customShort
	}

	val, _ := s.urlStorage.GetByID(ctx, id)
	if val != "" {
		return "", 0, ErrAlreadyInUse
	}

	if expiry == 0 {
		expiry = 24
	}

	err = s.urlStorage.SetByID(ctx, id, url, expiry*3600*time.Second)
	if err != nil {
		return "", 0, fmt.Errorf("cannot set URL by ID: %w")
	}

	remainingQuota, err := s.ipStorage.DecrAPIQuotaByIP(ctx, ip)
	if err != nil {
		return "", 0, fmt.Errorf("cannot decrement API quota by IP: %w", err)
	}

	return id, remainingQuota, nil
}
