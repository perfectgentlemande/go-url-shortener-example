package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

type Controller struct {
	srvc            service.Service
	UrlStorage      service.URLStorage
	IpStorage       service.IPStorage
	DefaultAPIQuota int
}

func (c *Controller) Ping(fCtx *fiber.Ctx) error {
	return fCtx.Status(fiber.StatusOK).JSON(map[string]string{
		"ping": "pong",
	})
}
