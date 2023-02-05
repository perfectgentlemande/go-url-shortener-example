package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

func (c *Controller) Resolve(fCtx *fiber.Ctx) error {
	ctx := fCtx.Context()
	url := fCtx.Params("url")

	value, err := c.srvc.Resolve(ctx, url)
	if err != nil {
		log.Printf("cannot resolve url: %s\n", err)
		if errors.Is(err, service.ErrNoSuchItem) {
			return fCtx.Status(fiber.StatusNotFound).JSON(APIError{Message: "short-url not found"})
		}

		return fCtx.Status(fiber.StatusInternalServerError).JSON(APIError{Message: "internal error"})
	}

	return fCtx.Redirect(value, http.StatusMovedPermanently)
}
