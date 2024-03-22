package server

import (
	"testing"

	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	el := &MockEventLoop{
		events: make(chan *loop.CacheEvent),
	}
	expectedAddr := ":8080"

	server := New(el, expectedAddr)

	assert.NotNil(t, server.eventLoop)
	assert.NotNil(t, server.httpServer)
	assert.Equal(t, expectedAddr, server.httpServer.Addr)
}
