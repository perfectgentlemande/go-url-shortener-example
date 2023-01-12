package routes

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/api/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/base62"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

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

func (c *Controller) Shorten(ctx *fiber.Ctx) error {
	dbCtx := context.TODO()
	body := &request{}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	limit, err := c.IpStorage.GetTTLByIP(dbCtx, ctx.IP())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get TTL"})
	}

	valInt, err := c.IpStorage.GetRequestsCountByIP(dbCtx, ctx.IP())
	if errors.Is(err, service.ErrNoSuchItem) {
		err = c.IpStorage.SetAPIQuotaByIP(dbCtx, ctx.IP(), c.DefaultAPIQuota, 30*60*time.Second)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot set api quota for IP"})
		}
	} else if err == nil {
		if valInt <= 0 {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Can't do that :)"})
	}

	// enforce HTTPS, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		id = base62.Encode(rand.Uint64())
	} else {
		id = body.CustomShort
	}

	val, _ := c.UrlStorage.GetByID(dbCtx, id)
	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL Custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = c.UrlStorage.SetByID(dbCtx, id, body.URL, body.Expiry*3600*time.Second)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot set URL by ID",
		})
	}

	remainingQuota, err := c.IpStorage.DecrAPIQuotaByIP(dbCtx, ctx.IP())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot decrAPI quota",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  int(remainingQuota),
		XRateLimitReset: limit / time.Nanosecond / time.Minute,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
