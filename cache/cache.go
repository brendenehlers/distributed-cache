package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"hash"
)

type Cache[K comparable, V any] interface {
	Read(key K) (V, error)
	Insert(key K, val V) error
	Remove(key K) error
}

type cacheEntry[K comparable, V any] struct {
	Key         K
	Val         V
	InitialHash uint32
	Deleted     bool
}

type InMemoryCache[K comparable, V any] struct {
	Size              int
	cache             []*cacheEntry[K, V]
	capacity          uint32
	resizeThreshold   float32
	resizeCoefficient uint16
	h                 hash.Hash
}

type Options struct {
	Capacity          uint32
	ResizeThreshold   float32
	ResizeCoefficient uint16
}

const DefaultCapacity = 1024
const DefaultResizeCoefficient = 2
const DefaultResizeThreshold = 0.75

func New[K comparable, V any](options Options) InMemoryCache[K, V] {
	if options.Capacity == 0 {
		options.Capacity = DefaultCapacity
	}
	if options.ResizeCoefficient == 0 {
		options.ResizeCoefficient = DefaultResizeCoefficient
	}
	if options.ResizeThreshold == 0.0 {
		options.ResizeThreshold = DefaultResizeThreshold
	}

	cache := InMemoryCache[K, V]{
		Size:              0,
		cache:             make([]*cacheEntry[K, V], options.Capacity),
		capacity:          options.Capacity,
		resizeThreshold:   options.ResizeThreshold,
		resizeCoefficient: options.ResizeCoefficient,
		h:                 sha256.New(),
	}

	return cache
}

func (c *InMemoryCache[K, V]) Read(key K) (V, error) {
	// get the starting index
	index, err := c.hash(key)
	if err != nil {
		var noop V
		return noop, err
	}

	// find the value with the matching key
	var x uint32 = 0
	for entry := c.cache[index]; entry.Key != key; entry = c.cache[index] {
		// iterated through the whole cache
		if x == c.capacity {
			panic("Cache capacity reached")
		}
		// found a nil entry before the key
		if entry == nil {
			panic("Element not found")
		}

		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	entry := c.cache[index]

	return entry.Val, nil
}

func (c *InMemoryCache[K, V]) Insert(key K, val V) error {
	// get the initial index
	index, err := c.hash(key)
	if err != nil {
		return err
	}

	// find the next open spot in the cache
	initialHashIndex := index
	var x uint32 = 0
	for entry := c.cache[index]; entry != nil && entry.Key != key; entry = c.cache[index] {
		// there's no space in the cache
		// this can happen if the resizeCoefficient is >= 1
		if x == c.capacity {
			panic("Cache capacity exceeded")
		}
		// find the next index
		x += 1
		index = (index + c.probing(x)) % c.capacity
		entry = c.cache[index]
	}

	entry := c.cache[index]
	if entry != nil {
		// update the existing entry
		entry.Val = val
		entry.Deleted = false
	} else {
		// insert the new entry
		c.cache[index] = &cacheEntry[K, V]{
			Key:         key,
			Val:         val,
			InitialHash: initialHashIndex,
			Deleted:     false,
		}
		c.Size += 1
	}

	return nil
}

func (c *InMemoryCache[K, V]) Remove(key K) error {
	panic("not implemented")
}

func (c *InMemoryCache[K, V]) encode(key K) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *InMemoryCache[K, V]) hash(key K) (uint32, error) {
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

func (c *InMemoryCache[K, V]) probing(x uint32) uint32 {
	// p(x) = x prevents propagation cycles
	return x
}
