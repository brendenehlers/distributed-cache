package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"hash"
	"time"
)

type Cache[K, V any] interface {
	Read(key K) (V, error)
	Insert(key K, val V) error
	Remove(key K) error
}

type cacheEntry[K, V any] struct {
	initialHash int
	key         K
	val         V
	lastUpdated time.Time
}

type InMemoryCache[K, V any] struct {
	cache             []cacheEntry[K, V]
	capacity          uint32
	size              uint32
	resizeThreshold   float32
	resizeCoefficient uint16
	h                 hash.Hash
}

type Options struct {
	Capacity          uint32
	ResizeThreshold   float32
	ResizeCoefficient uint16
}

func NewInMemoryCache[K, V any](options Options) InMemoryCache[K, V] {
	if options.Capacity == 0 {
		options.Capacity = 1024
	}
	if options.ResizeCoefficient == 0 {
		options.ResizeCoefficient = 2
	}
	if options.ResizeThreshold == 0.0 {
		options.ResizeThreshold = 0.67
	}

	cache := InMemoryCache[K, V]{
		cache:             make([]cacheEntry[K, V], options.Capacity),
		capacity:          options.Capacity,
		size:              0,
		resizeThreshold:   options.ResizeThreshold,
		resizeCoefficient: options.ResizeCoefficient,
		h:                 sha256.New(),
	}

	return cache
}

func (c InMemoryCache[K, V]) Read(key K) (V, error) {
	panic("not implemented")
}

func (c InMemoryCache[K, V]) Insert(key K, val V) error {
	panic("not implemented")
}

func (c InMemoryCache[K, V]) Remove(key K) error {
	panic("not implemented")
}

func (c InMemoryCache[K, V]) encode(key K) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c InMemoryCache[K, V]) hash(key K) (uint32, error) {
	encoded, err := c.encode(key)
	if err != nil {
		return 0, err
	}

	c.h.Write(encoded)
	hash := c.h.Sum(nil)
	index := binary.BigEndian.Uint32(hash) % c.capacity
	c.h.Reset()

	return index, nil
}

func (c InMemoryCache[K, V]) propagate(x uint32) uint32 {
	// p(x) = x prevents propagation cycles
	return x
}
