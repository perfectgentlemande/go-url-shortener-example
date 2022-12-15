package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/logger"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

func (c *Controller) resolve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	log := logger.GetLogger(ctx).WithField("id", id)

	ip, err := extractIP(r)
	if err != nil {
		log.WithError(err).Info("cannot extract ip")
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	urlStr, err := c.srvc.Resolve(ctx, id, ip)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.WithError(err).Info("no such slug")
			WriteError(ctx, w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		log.WithError(err).Info("cannot resolve URL by ID")
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, urlStr, http.StatusMovedPermanently)
}

func (c *Controller) shorten(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shorten"))
}
