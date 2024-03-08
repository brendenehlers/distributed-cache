package server

import (
	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
)

type EventLoop interface {
	Run()
	Send(event *loop.CacheEvent)
	Stop()
}
