package memory

import (
	"context"
	"testing"
	"time"

	goidempo "github.com/adelowo/go-idempo"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	duration := time.Hour
	cache, err := New(duration)

	require.NoError(t, err)
	require.NotNil(t, cache)

	memCache := cache.(*Memory)
	require.Equal(t, duration, memCache.duration)
	require.NotNil(t, memCache.items)
	require.Empty(t, memCache.items)
}

func TestMemory_Add(t *testing.T) {
	cache, err := New(time.Hour)
	require.NoError(t, err)

	item := goidempo.CacheItem{
		ID:           ulid.Make(),
		Key:          "test-key",
		FingerPrint:  "fp123",
		RequestBody:  "request body",
		ResponseBody: "response body",
		CreatedAt:    time.Now(),
	}

	err = cache.Add(context.Background(), item)
	require.NoError(t, err)

	retrieved, err := cache.Get(context.Background(), item.Key)
	require.NoError(t, err)

	require.Equal(t, item.ID, retrieved.ID)
	require.Equal(t, item.Key, retrieved.Key)
	require.Equal(t, item.FingerPrint, retrieved.FingerPrint)
	require.Equal(t, item.RequestBody, retrieved.RequestBody)
	require.Equal(t, item.ResponseBody, retrieved.ResponseBody)
}

func TestMemory_Get_NotFound(t *testing.T) {
	cache, err := New(time.Hour)
	require.NoError(t, err)

	_, err = cache.Get(context.Background(), "non-existent-key")
	require.Error(t, err)
	require.Equal(t, goidempo.ErrCacheKeyNotFound, err)
}

func TestMemory_Get_Success(t *testing.T) {
	cache, err := New(time.Hour)
	require.NoError(t, err)

	item := goidempo.CacheItem{
		ID:           ulid.Make(),
		Key:          "test-key",
		FingerPrint:  "fp123",
		RequestBody:  "request body",
		ResponseBody: "response body",
		CreatedAt:    time.Now(),
	}

	err = cache.Add(context.Background(), item)
	require.NoError(t, err)

	retrieved, err := cache.Get(context.Background(), item.Key)
	require.NoError(t, err)

	require.Equal(t, item.ID, retrieved.ID)
	require.Equal(t, item.Key, retrieved.Key)
	require.Equal(t, item.FingerPrint, retrieved.FingerPrint)
	require.Equal(t, item.RequestBody, retrieved.RequestBody)
	require.Equal(t, item.ResponseBody, retrieved.ResponseBody)
}

func TestMemory_Clear(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		cache, err := New(time.Hour)
		require.NoError(t, err)

		err = cache.Clear(context.Background())
		require.NoError(t, err)

		memCache := cache.(*Memory)
		require.Empty(t, memCache.items)
	})

	t.Run("no expired items", func(t *testing.T) {
		cache, err := New(time.Hour)
		require.NoError(t, err)

		item := goidempo.CacheItem{
			ID:           ulid.Make(),
			Key:          "test-key",
			FingerPrint:  "fp123",
			RequestBody:  "request body",
			ResponseBody: "response body",
			CreatedAt:    time.Now(),
		}

		err = cache.Add(context.Background(), item)
		require.NoError(t, err)

		err = cache.Clear(context.Background())
		require.NoError(t, err)

		retrieved, err := cache.Get(context.Background(), item.Key)
		require.NoError(t, err)
		require.Equal(t, item.ID, retrieved.ID)
	})

	t.Run("with expired items", func(t *testing.T) {
		duration := 10 * time.Millisecond
		cache, err := New(duration)
		require.NoError(t, err)

		item1 := goidempo.CacheItem{
			ID:           ulid.Make(),
			Key:          "expired-key",
			FingerPrint:  "fp123",
			RequestBody:  "request body",
			ResponseBody: "response body",
			CreatedAt:    time.Now().Add(-time.Hour),
		}

		item2 := goidempo.CacheItem{
			ID:           ulid.Make(),
			Key:          "fresh-key",
			FingerPrint:  "fp456",
			RequestBody:  "request body 2",
			ResponseBody: "response body 2",
			CreatedAt:    time.Now(),
		}

		err = cache.Add(context.Background(), item1)
		require.NoError(t, err)
		err = cache.Add(context.Background(), item2)
		require.NoError(t, err)

		time.Sleep(duration + time.Millisecond)

		err = cache.Clear(context.Background())
		require.NoError(t, err)

		_, err = cache.Get(context.Background(), item1.Key)
		require.Error(t, err)
		require.Equal(t, goidempo.ErrCacheKeyNotFound, err)

		_, err = cache.Get(context.Background(), item2.Key)
		require.Error(t, err)
		require.Equal(t, goidempo.ErrCacheKeyNotFound, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		cache, err := New(time.Hour)
		require.NoError(t, err)

		for i := range 100 {
			item := goidempo.CacheItem{
				ID:           ulid.Make(),
				Key:          goidempo.IdempotencyKey(string(rune(i))),
				FingerPrint:  "fp123",
				RequestBody:  "request body",
				ResponseBody: "response body",
				CreatedAt:    time.Now().Add(-time.Hour),
			}
			err = cache.Add(context.Background(), item)
			require.NoError(t, err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = cache.Clear(ctx)
		require.Error(t, err)
		require.Equal(t, context.Canceled, err)
	})
}

func TestMemory_Add_KeyConflict(t *testing.T) {
	cache, err := New(time.Hour)
	require.NoError(t, err)

	key := goidempo.IdempotencyKey("test-key")

	item1 := goidempo.CacheItem{
		ID:           ulid.Make(),
		Key:          key,
		FingerPrint:  "fp123",
		RequestBody:  "request body 1",
		ResponseBody: "response body 1",
		CreatedAt:    time.Now(),
	}

	item2 := goidempo.CacheItem{
		ID:           ulid.Make(),
		Key:          key,
		FingerPrint:  "fp456",
		RequestBody:  "request body 2",
		ResponseBody: "response body 2",
		CreatedAt:    time.Now(),
	}

	err = cache.Add(context.Background(), item1)
	require.NoError(t, err)

	err = cache.Add(context.Background(), item2)
	require.Error(t, err)
	require.Equal(t, goidempo.ErrKeyConflict, err)

	retrieved, err := cache.Get(context.Background(), key)
	require.NoError(t, err)
	require.Equal(t, item1.ID, retrieved.ID)
	require.Equal(t, item1.FingerPrint, retrieved.FingerPrint)
}
