package api

import (
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

type Controller struct {
	srvc *service.Service
}

func New(srvc *service.Service) *Controller {
	return &Controller{
		srvc: srvc,
	}
}
