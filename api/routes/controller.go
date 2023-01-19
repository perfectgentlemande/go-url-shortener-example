package routes

import (
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"
)

type Controller struct {
	srvc            service.Service
	UrlStorage      service.URLStorage
	IpStorage       service.IPStorage
	DefaultAPIQuota int
}
