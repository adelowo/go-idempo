package goidempo

import (
	"net/http"
	"time"
)

const (
	defaultExpiration = time.Hour * 24
)

type SkipFn func(*http.Request) bool

var (
	DefaultSkipFn = func(r *http.Request) bool { return false }
)

type handler struct {
	expirationTime time.Duration
	skipFn         SkipFn
}

func Idempotency(opts ...Option) func(next http.Handler) http.Handler {
	h := &handler{
		expirationTime: defaultExpiration,
		skipFn:         DefaultSkipFn,
	}

	for _, opt := range opts {
		opt(h)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if h.skipFn(r) {
				next.ServeHTTP(w, r)
				return
			}

			_, _ = KeyFromRequest(r.Header)

		})
	}
}
