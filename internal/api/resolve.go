package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

func (c *Controller) Resolve(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	value, err := c.srvc.Resolve(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchItem) {
			log.Printf("short-url not found\n")
			WriteError(ctx, w, http.StatusNotFound, "short-url not found")
			return
		}

		log.Printf("cannot resolve url: %s\n", err)
		WriteError(ctx, w, http.StatusInternalServerError, "short-url not found")
		return
	}

	http.Redirect(w, r, value, http.StatusMovedPermanently)
}
