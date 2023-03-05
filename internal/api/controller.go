package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

type Config struct {
	AppPort string
	Domain  string
}

type Controller struct {
	srvc   *service.Service
	domain string
}

func New(srvc *service.Service, domain string) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
func WriteError(ctx context.Context, w http.ResponseWriter, status int, message string) {
	err := RespondWithJSON(w, status, APIError{Message: message})
	if err != nil {
		log.Printf("cannot write response: %s\n", err)
	}
}
func WriteSuccessful(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	err := RespondWithJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("cannot write response: %s\n", err)
	}
}
