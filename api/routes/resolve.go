package routes

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) Resolve(ctx *fiber.Ctx) error {
	dbCtx := context.TODO()
	url := ctx.Params("url")

	value, err := c.R.Get(dbCtx, url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short-url not found in db"})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
	}

	_ = c.RInr.Incr(dbCtx, "counter")

	return ctx.Redirect(value, http.StatusMovedPermanently)
}
