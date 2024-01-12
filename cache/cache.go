package cache

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"sync"
)

type Cache[K comparable, V any] interface {
	Read(key K) (V, error)
	Insert(key K, val V) error
	Remove(key K) error
}

type cacheEntry[K comparable, V any] struct {
	Key              K
	Val              V
	InitialHashIndex uint32
	Deleted          bool
}

type InMemoryCache[K comparable, V any] struct {
	Size              int
	cache             []*cacheEntry[K, V]
	capacity          uint32
	resizeThreshold   float32
	resizeCoefficient uint32
	mux               sync.RWMutex
}

type Options struct {
	Capacity          uint32
	ResizeThreshold   float32
	ResizeCoefficient uint32
}

const DefaultCapacity = 1024
const DefaultResizeCoefficient = 2
const DefaultResizeThreshold = 0.75

func New[K comparable, V any](options Options) *InMemoryCache[K, V] {
	if options.Capacity == 0 {
		options.Capacity = DefaultCapacity
	}
	if options.ResizeCoefficient == 0 {
		options.ResizeCoefficient = DefaultResizeCoefficient
	}
	if options.ResizeThreshold == 0.0 {
		options.ResizeThreshold = DefaultResizeThreshold
	}

	cache := &InMemoryCache[K, V]{
		Size:              0,
		cache:             make([]*cacheEntry[K, V], options.Capacity),
		capacity:          options.Capacity,
		resizeThreshold:   options.ResizeThreshold,
		resizeCoefficient: options.ResizeCoefficient,
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
	var x uint32 = 1
	c.mux.RLock()
	for entry := c.cache[index]; entry.Key != key; entry = c.cache[index] {
		// iterated through the whole cache
		if x == c.capacity {
			panic("Cache capacity reached")
		}
		// found a nil entry before the key
		if entry == nil {
			var noop V
			return noop, nil
		}

		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	entry := c.cache[index]
	c.mux.RUnlock()

	return entry.Val, nil
}

func (c *InMemoryCache[K, V]) Insert(key K, val V) error {
	if c.Size/int(c.capacity) >= int(c.resizeThreshold) {
		c.increaseCacheSize()
	}

	// get the initial index
	index, err := c.hash(key)
	if err != nil {
		return err
	}

	// find the next open spot in the cache
	initialHashIndex := index
	var x uint32 = 1
	c.mux.RLock()
	for entry := c.cache[index]; entry != nil && entry.Key != key && !entry.Deleted; entry = c.cache[index] {
		// there's no space in the cache
		// this can happen if the resizeCoefficient is >= 1
		if x == c.capacity {
			panic("Cache capacity exceeded")
		}
		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	entry := c.cache[index]
	c.mux.RUnlock()

	if entry != nil {
		c.mux.Lock()
		// update the existing entry
		entry.Val = val
		entry.Deleted = false
		c.mux.Unlock()
	} else {
		c.mux.Lock()
		// insert the new entry
		c.cache[index] = &cacheEntry[K, V]{
			Key:              key,
			Val:              val,
			InitialHashIndex: initialHashIndex,
			Deleted:          false,
		}
		c.Size += 1
		c.mux.Unlock()
	}

	return nil
}

func (c *InMemoryCache[K, V]) Remove(key K) error {
	// get the initial index
	index, err := c.hash(key)
	if err != nil {
		return err
	}

	// find the entry in the cache
	var x uint32 = 1
	c.mux.RLock()
	for entry := c.cache[index]; entry != nil && entry.Key != key; entry = c.cache[index] {
		// there's no space in the cache
		// this can happen if the resizeCoefficient is >= 1
		if x == c.capacity {
			panic("Cache capacity exceeded")
		}
		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	// mark the entry as deleted
	entry := c.cache[index]
	c.mux.RUnlock()

	if entry != nil {
		c.mux.Lock()
		entry.Deleted = true
		c.Size -= 1
		c.mux.Unlock()
	}

	return nil
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

	h := sha256.New()
	h.Write(encoded)
	hash := h.Sum(nil)
	index := binary.BigEndian.Uint32(hash) % c.capacity
	h.Reset() // don't know if this is needed

	return index, nil
}

func (c *InMemoryCache[K, V]) probing(x uint32) uint32 {
	// p(x) = x prevents propagation cycles
	return x
}

func (c *InMemoryCache[K, V]) increaseCacheSize() {
	// create a new cache with the increased size
	newCache := make([]*cacheEntry[K, V], c.resizeCoefficient*c.capacity)

	c.mux.RLock()
	// add the old values to the new cache
	for _, oldCacheEntry := range c.cache {
		// skip nil and deleted entries
		if oldCacheEntry == nil || oldCacheEntry.Deleted {
			continue
		}

		index := oldCacheEntry.InitialHashIndex

		// find the location in the new cache
		var x uint32 = 1
		for entry := newCache[index]; entry != nil; entry = newCache[index] {
			index = (index + c.probing(x)) % c.capacity
			x += 1
		}

		newCache[index] = oldCacheEntry
	}
	c.mux.RUnlock()

	c.mux.Lock()
	// update the old cache to the new cache
	c.cache = newCache
	c.mux.Unlock()
}
