package main

import (
	"github.com/brendenehlers/go-distributed-cache/cache-node/data"
	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
	"github.com/brendenehlers/go-distributed-cache/cache-node/server"
)

// no reason for this other than I wanted to practice the adapter pattern
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

func main() {

	inMemoryCache := data.NewInMemoryCache[string, loop.CacheEntry](data.Options{})
	cache := NewInMemoryCacheAdapter(inMemoryCache)
	eventLoop := loop.NewEventLoop(cache)
	server := server.NewServer(eventLoop, ":8080")

	server.Run()
}
