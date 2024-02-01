package main

import (
	"fmt"

	"github.com/brendenehlers/go-distributed-cache/cache-node/cache"
)

func main() {
	cache := cache.New[int, string](cache.Options{
		Capacity: 4,
	})

	kvs := make(map[int]string)
	kvs[123897] = "hello world 1"
	kvs[981273] = "hello world 2"
	kvs[871263] = "hello world 3"
	kvs[182963] = "hello world 4"

	for k, v := range kvs {
		if err := cache.Insert(k, v); err != nil {
			println(err)
			return
		}
	}

	for k := range kvs {
		val, err := cache.Read(k)
		if err != nil {
			println(err)
			return
		}
		fmt.Printf("The value for '%d' is '%s'\n", k, val)
	}

}
