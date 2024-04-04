package adapter

import (
	"github.com/brendenehlers/go-distributed-cache/cache-node/data"
	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
)

type InMemoryCacheAdapter struct {
	inMemoryCache *data.InMemoryCache[string, loop.CacheEntry]
}

func (adapter *InMemoryCacheAdapter) Get(key string) (loop.CacheEntry, bool) {
	return adapter.inMemoryCache.Read(key)
}

func (adapter *InMemoryCacheAdapter) Set(key string, val loop.CacheEntry) error {
	return adapter.inMemoryCache.Insert(key, val)
}

func (adapter *InMemoryCacheAdapter) Delete(key string) error {
	return adapter.inMemoryCache.Remove(key)
}

func NewInMemoryCacheAdapter(cache *data.InMemoryCache[string, loop.CacheEntry]) loop.Cache {
	return &InMemoryCacheAdapter{
		inMemoryCache: cache,
	}
}