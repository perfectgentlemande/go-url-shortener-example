package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/logger"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"github.com/sirupsen/logrus"
)

type Params struct {
	Addr string
	Log  *logrus.Entry
	Srvc *service.Service
}

type Controller struct {
	srvc *service.Service
}

func NewController(srvc *service.Service) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func NewServer(params *Params) *http.Server {
	ctrl := NewController(params.Srvc)

	apiRouter := chi.NewRouter()
	apiRouter.Use(logger.NewLoggingMiddleware(params.Log))
	apiRouter.Get("/{id}", ctrl.resolve)
	apiRouter.Post("/api/v1", ctrl.shorten)

	return &http.Server{
		Addr:    params.Addr,
		Handler: apiRouter,
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
func WriteError(ctx context.Context, w http.ResponseWriter, status int, message string) {
	log := logger.GetLogger(ctx)

	err := RespondWithJSON(w, status, APIError{Message: message})
	if err != nil {
		log.WithError(err).Error("write response error")
	}
}
func WriteSuccessful(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	log := logger.GetLogger(ctx)

	err := RespondWithJSON(w, http.StatusOK, payload)
	if err != nil {
		log.WithError(err).Error("write response error")
	}
}
