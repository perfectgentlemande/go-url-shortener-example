package routes

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/api/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/base62"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// APIError defines model for APIError.
type APIError struct {
	// Error message
	Message string `json:"message"`
}

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func (c *Controller) Shorten(fCtx *fiber.Ctx) error {
	ctx := fCtx.Context()
	body := &request{}

	if err := fCtx.BodyParser(&body); err != nil {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "cannot parse JSON"})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "Invalid URL"})
	}

	limit, err := c.IpStorage.GetTTLByIP(ctx, fCtx.IP())
	if err != nil {
		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "cannot get TTL"})
	}

	valInt, err := c.IpStorage.GetRequestsCountByIP(ctx, fCtx.IP())
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

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  int(remainingQuota),
		XRateLimitReset: limit / time.Nanosecond / time.Minute,
	}

	return fCtx.Status(fiber.StatusOK).JSON(resp)
}
