package main

import (
	"fmt"

	"github.com/brendenehlers/go-distributed-cache/cache-node/data"
	"github.com/brendenehlers/go-distributed-cache/cache-node/domain"
)

// no reason for this other than I wanted to practice the adapter pattern
type InMemoryCacheAdapter struct {
	inMemoryCache *data.InMemoryCache[string, domain.CacheEntry]
}

func (adapter *InMemoryCacheAdapter) Get(key string) (domain.CacheEntry, bool) {
	return adapter.inMemoryCache.Read(key)
}

func (adapter *InMemoryCacheAdapter) Set(key string, val domain.CacheEntry) error {
	return adapter.inMemoryCache.Insert(key, val)
}

func (adapter *InMemoryCacheAdapter) Delete(key string) error {
	return adapter.inMemoryCache.Remove(key)
}

func NewInMemoryCacheAdapter(cache *data.InMemoryCache[string, domain.CacheEntry]) domain.Cache {
	return &InMemoryCacheAdapter{
		inMemoryCache: cache,
	}
}

func main() {

	inMemoryCache := data.NewInMemoryCache[string, domain.CacheEntry](data.Options{})
	cache := NewInMemoryCacheAdapter(inMemoryCache)
	loop := domain.NewEventLoop(cache)

	go loop.Run()

	responseChan := make(chan domain.EventResponse)
	errorChan := make(chan error)

	loop.Send(domain.Event{
		Type:      "set",
		Key:       "hello",
		Val:       "world",
		ErrorChan: errorChan,
	})

	loop.Send(domain.Event{
		Type:         "get",
		Key:          "hello",
		ResponseChan: responseChan,
		ErrorChan:    errorChan,
	})

	if resp := <-responseChan; !resp.Ok {
		fmt.Println("error occurred when fetching value")
	} else {
		fmt.Printf("Cache value: %s", resp.Value)
	}
}
