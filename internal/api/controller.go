package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

type Controller struct {
	srvc *service.Service
}

func New(srvc *service.Service) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func (c *Controller) Ping(fCtx *fiber.Ctx) error {
	return fCtx.Status(fiber.StatusOK).JSON(map[string]string{
		"ping": "pong",
	})
}
