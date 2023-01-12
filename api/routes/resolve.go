package routes

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"
)

func (c *Controller) Resolve(fCtx *fiber.Ctx) error {
	ctx := fCtx.Context()
	url := fCtx.Params("url")

	value, err := c.UrlStorage.GetByID(ctx, url)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchItem) {
			return fCtx.Status(fiber.StatusNotFound).JSON(APIError{Message: "short-url not found in db"})
		}

		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "Internal error"})
	}

	err = c.IpStorage.IncrRequestCounter(ctx)
	if err != nil {
		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "Cannot increment request counter"})
	}

	return fCtx.Redirect(value, http.StatusMovedPermanently)
}
