package main

import (
	"fmt"

	"github.com/brendenehlers/go-distributed-cache/cache-node/data"
	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
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

	go eventLoop.Run()

	event1, resChan1, errChan1 := loop.CreateSetEvent("hello", "world")
	eventLoop.Send(event1)

	select {
	case resp := <-resChan1:
		fmt.Printf("set ok: %v, val: %v\n", resp.Ok, resp.Value)
	case err := <-errChan1:
		panic(err)
	}

	event2, resChan2, errChan2 := loop.CreateGetEvent("hello")
	eventLoop.Send(event2)

	select {
	case resp := <-resChan2:
		fmt.Printf("get ok: %v, val: %v\n", resp.Ok, resp.Value)
	case err := <-errChan2:
		panic(err)
	}

	// server := presentation.NewServer(loop, ":8080")
	// server.StartServerAndLoop()
}
