package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber"
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

func (s *Service) Shorten(ctx context.Context, ip, url string) (string, error) {
	limit, err := s.ipStorage.GetTTLByIP(ctx, ip)
	if err != nil {
		return "", fmt.Errorf("cannot get TTL by IP: %w", err)
	}

	valInt, err := s.ipStorage.GetRequestsCountByIP(ctx, ip)
	if errors.Is(err, service.ErrNoSuchItem) {
		err = c.IpStorage.SetAPIQuotaByIP(ctx, fCtx.IP(), c.DefaultAPIQuota, 30*60*time.Second)
		if err != nil {
			return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "cannot set api quota for IP"})
		}
	} else if err == nil {
		if valInt <= 0 {
			return fCtx.Status(fiber.StatusServiceUnavailable).JSON(
				APIError{Message: fmt.Sprintf("Rate limit exceeded, rate_limit_reset: %d", limit/time.Nanosecond/time.Minute)})
		}
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return fCtx.Status(fiber.StatusServiceUnavailable).JSON(APIError{Message: "Can't do that :)"})
	}

	// enforce HTTPS, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		id = base62.Encode(rand.Uint64())
	} else {
		id = body.CustomShort
	}

	val, _ := c.UrlStorage.GetByID(ctx, id)
	if val != "" {
		return fCtx.Status(fiber.StatusForbidden).JSON(APIError{Message: "URL Custom short is already in use"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = c.UrlStorage.SetByID(ctx, id, body.URL, body.Expiry*3600*time.Second)
	if err != nil {
		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "Cannot set URL by ID"})
	}

	remainingQuota, err := c.IpStorage.DecrAPIQuotaByIP(ctx, fCtx.IP())
	if err != nil {
		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "Cannot decrAPI quota"})
	}

}
