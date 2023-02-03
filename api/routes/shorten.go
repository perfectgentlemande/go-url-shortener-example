package routes

import (
	"errors"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/api/helpers"
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

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "Can't do that :)"})
	}

	newID, remainingQuota, limit, err := c.srvc.Shorten(ctx, fCtx.IP(), body.URL, body.CustomShort, body.Expiry)
	if err != nil {
		if errors.Is(err, service.ErrRateLimitExceeded) {
			return fCtx.Status(fiber.StatusTooManyRequests).JSON(APIError{Message: "rate limit exceeded"})
		}
		if errors.Is(err, service.ErrAlreadyInUse) {
			return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "already in use"})
		}

		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "cannot shorten URL"})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + newID,
		Expiry:          body.Expiry,
		XRateRemaining:  int(remainingQuota),
		XRateLimitReset: limit / time.Nanosecond / time.Minute,
	}

	return fCtx.Status(fiber.StatusOK).JSON(resp)
}
