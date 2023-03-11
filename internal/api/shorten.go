package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"

	"github.com/asaskevich/govalidator"
)

func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := GetLogger(ctx)
	shortenReq := &ShortenRequest{}

	if err := json.NewDecoder(r.Body).Decode(&shortenReq); err != nil {
		log.WithError(err).Debug("wrong user data")
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	log = log.WithField("url", shortenReq.Url)

	// check if the input is an actual URL
	if !govalidator.IsURL(shortenReq.Url) {
		log.Debug("invalid URL")
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// check for domain error
	if !helpers.RemoveDomainError(shortenReq.Url, c.domain) {
		log.Debug("can't do that :)")
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newID, remainingQuota, limit, err := c.srvc.Shorten(ctx, r.RemoteAddr, shortenReq.Url, shortenReq.Short, time.Duration(shortenReq.Expiry))
	if err != nil {

		if errors.Is(err, service.ErrRateLimitExceeded) {
			log.WithError(err).Debug("rate limit exceeded")
			WriteError(ctx, w, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			return
		}
		if errors.Is(err, service.ErrAlreadyInUse) {
			log.WithError(err).Debug("already in use")
			WriteError(ctx, w, http.StatusBadRequest, "already in use")
			return
		}

		log.WithError(err).Error("cannot shorten URL")
		WriteError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	WriteSuccessful(ctx, w, ShortenResponse{
		ShortenRequest: ShortenRequest{
			Url:    shortenReq.Url,
			Short:  os.Getenv("DOMAIN") + "/" + newID,
			Expiry: shortenReq.Expiry,
		},
		RateLimitRemaining: int64(remainingQuota),
		RateLimitReset:     int64(limit / time.Nanosecond / time.Minute),
	})
}
