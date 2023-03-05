package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/helpers"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"

	"github.com/asaskevich/govalidator"
)

func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortenReq := &ShortenRequest{}

	if err := json.NewDecoder(r.Body).Decode(&shortenReq); err != nil {
		log.Printf("wrong user data: %s\n", err)
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(shortenReq.Url) {
		log.Printf("invalid URL\n")
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// check for domain error
	if !helpers.RemoveDomainError(shortenReq.Url, c.domain) {
		log.Printf("can't do that :)\n")
		WriteError(ctx, w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newID, remainingQuota, limit, err := c.srvc.Shorten(ctx, r.RemoteAddr, shortenReq.Url, shortenReq.Short, time.Duration(shortenReq.Expiry))
	if err != nil {
		log.Printf("cannot shorten URL: %s\n", err)
		if errors.Is(err, service.ErrRateLimitExceeded) {
			log.Printf("rate limit exceeded\n")
			WriteError(ctx, w, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			return
		}
		if errors.Is(err, service.ErrAlreadyInUse) {
			log.Printf("already in use\n")
			WriteError(ctx, w, http.StatusBadRequest, "already in use")
			return
		}

		log.Printf("cannot shorten URL: %s\n", err)
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
