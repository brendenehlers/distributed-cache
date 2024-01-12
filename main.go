package main

import (
	"github.com/brendenehlers/go-distributed-cache/cache"
)

func main() {
	cache := cache.NewInMemoryCache[string, string](cache.Options{})
}
