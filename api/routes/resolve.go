package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"
)

func (c *Controller) Resolve(ctx *fiber.Ctx) error {
	dbCtx := context.TODO()
	url := ctx.Params("url")

	value, err := c.UrlStorage.GetByID(dbCtx, url)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchItem) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short-url not found in db"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
	}

	err = c.IpStorage.IncrRequestCounter(dbCtx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot increment request counter"})
	}

	return ctx.Redirect(value, http.StatusMovedPermanently)
}
