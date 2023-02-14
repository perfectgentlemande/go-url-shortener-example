package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortenReq := &ShortenRequest{}

	if err := fCtx.BodyParser(&shortenReq); err != nil {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "cannot parse JSON"})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(shortenReq.Url) {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "Invalid URL"})
	}

	// check for domain error
	if !helpers.RemoveDomainError(shortenReq.Url) {
		return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "Can't do that :)"})
	}

	newID, remainingQuota, limit, err := c.srvc.Shorten(ctx, fCtx.IP(), shortenReq.Url, shortenReq.Short, shortenReq.Expiry)
	if err != nil {
		log.Printf("cannot shorten URL: %s\n", err)
		if errors.Is(err, service.ErrRateLimitExceeded) {
			return fCtx.Status(fiber.StatusTooManyRequests).JSON(APIError{Message: "rate limit exceeded"})
		}
		if errors.Is(err, service.ErrAlreadyInUse) {
			return fCtx.Status(fiber.StatusBadRequest).JSON(APIError{Message: "already in use"})
		}

		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "cannot shorten URL"})
	}

	resp := ShortenResponse{
		ShortenRequest: ShortenRequest{
			Url:    shortenReq.Url,
			Short:  os.Getenv("DOMAIN") + "/" + newID,
			Expiry: shortenReq.Expiry,
		},
		RateLimitRemaining: int64(remainingQuota),
		RateLimitReset:     int64(limit / time.Nanosecond / time.Minute),
	}

	return fCtx.Status(fiber.StatusOK).JSON(resp)
}
