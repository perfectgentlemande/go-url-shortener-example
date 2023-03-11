package api

import (
	"errors"
	"net/http"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

func (c *Controller) Resolve(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	log := GetLogger(ctx)

	value, err := c.srvc.Resolve(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchItem) {
			log.WithError(err).Debug("short-url not found")
			WriteError(ctx, w, http.StatusNotFound, "short-url not found")
			return
		}

		log.WithError(err).Error("cannot resolve url")
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, value, http.StatusMovedPermanently)
}
