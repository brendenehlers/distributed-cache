package main

import (
	"flag"
	"fmt"

	"github.com/brendenehlers/go-distributed-cache/cache-node/adapter"
	"github.com/brendenehlers/go-distributed-cache/cache-node/data"
	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
	"github.com/brendenehlers/go-distributed-cache/cache-node/server"
)

var (
	portFlag = flag.Int("port", 8080, "port for server to listen on")
	hostnameFlag = flag.String("hostname", "localhost", "hostname for the server")
)

func main() {	
	flag.Parse()

	inMemoryCache := data.NewInMemoryCache[string, loop.CacheEntry](data.Options{})
	cache := adapter.NewInMemoryCacheAdapter(inMemoryCache)
	eventLoop := loop.NewEventLoop(cache)

	host := fmt.Sprintf("%s:%d", *hostnameFlag, *portFlag)
	server := server.New(eventLoop, host)

	server.Run()
}
