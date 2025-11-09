package goidempo

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

type CacheItem struct {
	ID           ulid.ULID      `json:"id"`
	Key          IdempotencyKey `json:"key"`
	FingerPrint  string         `json:"finger_print"`
	RequestBody  string         `json:"request_body"`
	ResponseBody string         `json:"response_body"`
	CreatedAt    time.Time      `json:"created_at"`
}

// ENUM(database,redis,memory)
type CacheProvider string

type Cache interface {
	Add(context.Context, CacheItem) error
	Get(context.Context, IdempotencyKey) (CacheItem, error)
}
