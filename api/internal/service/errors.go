package service

import "errors"

var (
	ErrAlreadyInUse      = errors.New("custom short URL is already in use")
	ErrNoSuchItem        = errors.New("no such item")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
