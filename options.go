package goidempo

import "time"

type Option func(*handler)

// WithExpiration defines how long a cached value should exists for
// and controls purging of the data
// https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-02#section-2.3
func WithExpiration(t time.Duration) Option {
	return func(h *handler) {
		h.expirationTime = t
	}
}

// WithSkipFn allows you configure when idempotency checks
// should be skipped.
//
// By default no requests are skipped but in theory, you
// could filter bby path and what not
func WithSkipFn(s SkipFn) Option {
	return func(h *handler) {
		h.skipFn = s
	}
}
