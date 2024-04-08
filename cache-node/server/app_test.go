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
	expectedReg := "asdf"

	server := New(el, expectedAddr, expectedReg)

	assert.NotNil(t, server.eventLoop)
	assert.NotNil(t, server.Server)
	assert.Equal(t, expectedAddr, server.Server.Addr)
	assert.Equal(t, expectedReg, server.registryUrl)
}
