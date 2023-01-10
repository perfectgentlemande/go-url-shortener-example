package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/perfectgentlemande/go-url-shortener-example/api/internal/service"
)

type Controller struct {
	UrlStorage service.URLStorage
	R          *redis.Client
	RInr       *redis.Client
	R2         *redis.Client
}
