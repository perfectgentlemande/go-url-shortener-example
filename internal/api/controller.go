package api

import (
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
)

type Controller struct {
	srvc            service.Service
	UrlStorage      service.URLStorage
	IpStorage       service.IPStorage
	DefaultAPIQuota int
}
