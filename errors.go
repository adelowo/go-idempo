package goidempo

import "errors"

type IdempotencyError error

var (
	ErrKeyNotFound = errors.New("no Idempotency key found in request")
)
