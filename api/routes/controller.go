package routes

import (
	"github.com/go-redis/redis/v8"
)

type Controller struct {
	urlStorage service.URLStorage
	R          *redis.Client
	RInr       *redis.Client
	R2         *redis.Client
}
