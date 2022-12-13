package api

import (
	"errors"
	"net/http"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/logger"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

func (c *Controller) resolve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	urlStr := r.URL.String()

	urlStr, err := c.srvc.GetURLByID(ctx, urlStr)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.WithError(err).Info("no such slug")
			WriteError(ctx, w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		log.WithError(err).Info("cannot get URL by ID")
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, urlStr, http.StatusMovedPermanently)
}

// func Resolve(ctx *fiber.Ctx) error {
// 	url := ctx.Params("url")

// 	r := database.CreateClient(0)
// 	defer r.Close()

// 	value, err := r.Get(database.Ctx, url).Result()
// 	if err == redis.Nil {
// 		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short-url not found in db"})
// 	} else if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
// 	}

// 	rInr := database.CreateClient(1)
// 	defer rInr.Close()

// 	_ = rInr.Incr(database.Ctx, "counter")

// 	return ctx.Redirect(value, 301)
// }

func (c *Controller) shorten(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shorten"))
}
