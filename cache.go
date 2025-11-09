package goidempo

import "context"

// ENUM(database,redis,memory)
type CacheProvider string

type Cache interface {
	Add(context.Context) error
}
