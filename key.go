package goidempo

import (
	"net/http"

	"github.com/ayinke-llc/hermes"
)

const headerName = "Idempotency-Key"

func KeyFromRequest(h http.Header) (IdempotencyKey, error) {

	key := h.Get(headerName)
	if hermes.IsStringEmpty(key) {
		return "", ErrKeyNotFound
	}

	return IdempotencyKey(key), nil
}

type IdempotencyKey string
