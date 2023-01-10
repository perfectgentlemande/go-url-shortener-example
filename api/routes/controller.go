package routes

import "github.com/go-redis/redis/v8"

type Controller struct {
	r    *redis.Client
	rInr *redis.Client
	r1   *redis.Client
	r2   *redis.Client
}
