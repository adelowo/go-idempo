package memory

import (
	"context"
	"sync"
	"time"

	goidempo "github.com/adelowo/go-idempo"
)

type Memory struct {
	mutex    sync.RWMutex
	items    map[goidempo.IdempotencyKey]goidempo.CacheItem
	duration time.Duration
}

func New(t time.Duration) (goidempo.Cache, error) {
	return &Memory{
		mutex:    sync.RWMutex{},
		items:    make(map[goidempo.IdempotencyKey]goidempo.CacheItem),
		duration: t,
	}, nil
}

func (m *Memory) Add(_ context.Context, item goidempo.CacheItem) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.items[item.Key]; exists {
		return goidempo.ErrKeyConflict
	}

	m.items[item.Key] = item
	return nil
}

func (m *Memory) Get(_ context.Context, key goidempo.IdempotencyKey) (goidempo.CacheItem, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return goidempo.CacheItem{}, goidempo.ErrCacheKeyNotFound
	}

	return item, nil
}

func (m *Memory) Clear(ctx context.Context) error {
	var deleteKeys []goidempo.IdempotencyKey

	m.mutex.RLock()

	for k, v := range m.items {
		select {
		case <-ctx.Done():
			m.mutex.RUnlock()
			return ctx.Err()
		default:
			if time.Since(v.CreatedAt) > m.duration {
				deleteKeys = append(deleteKeys, k)
			}
		}
	}

	m.mutex.RUnlock()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(deleteKeys) > 0 {
		for _, key := range deleteKeys {
			delete(m.items, key)
		}
	}

	return nil
}
