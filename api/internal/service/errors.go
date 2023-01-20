package service

import "errors"

var (
	ErrAlreadyInUse      = errors.New("custom short URL is already in use")
	ErrCantDoThat        = errors.New("can't do that :)")
	ErrNoSuchItem        = errors.New("no such item")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
