package api

import (
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
	apiRouter.Get("/:url", ctrl.resolve)
	apiRouter.Post("/api/v1", ctrl.shorten)

	return &http.Server{
		Addr:    params.Addr,
		Handler: apiRouter,
	}
}
