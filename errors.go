package goidempo

import "errors"

type IdempotencyError error

var (
	ErrKeyNotFound      = errors.New("no Idempotency key found in request")
	ErrCacheKeyNotFound = errors.New("idempotency key not previously stored")
	ErrKeyConflict      = errors.New("idempotency key exists ")
)
