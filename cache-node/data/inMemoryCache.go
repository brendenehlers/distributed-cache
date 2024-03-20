package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"sync"
)

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

func NewInMemoryCache[K comparable, V any](options Options) *InMemoryCache[K, V] {
	options = assignDefaultOptions[K, V](options)
	cache := buildCache[K, V](options)

	return cache
}

func assignDefaultOptions[K comparable, V any](options Options) Options {
	if options.Capacity == 0 {
		options.Capacity = DefaultCapacity
	}
	if options.ResizeCoefficient == 0 {
		options.ResizeCoefficient = DefaultResizeCoefficient
	}
	if options.ResizeThreshold == 0.0 {
		options.ResizeThreshold = DefaultResizeThreshold
	}

	return options
}

func buildCache[K comparable, V any](options Options) *InMemoryCache[K, V] {
	return &InMemoryCache[K, V]{
		Size:              0,
		cache:             make([]*cacheEntry[K, V], options.Capacity),
		capacity:          options.Capacity,
		resizeThreshold:   options.ResizeThreshold,
		resizeCoefficient: options.ResizeCoefficient,
	}
}

func (c *InMemoryCache[K, V]) Insert(key K, val V) error {
	c.checkCacheSize()

	// get the initial index
	index, err := c.hash(key)
	if err != nil {
		return err
	}

	// find the next open spot in the cache
	initialHashIndex := index
	index = c.findNextEmptySpotInCache(key, index)
	entry := c.cache[index]

	if entry != nil {
		c.updateCacheEntry(entry, val)
	} else {
		c.insertNewCacheEntry(index, c.createNewCacheEntry(key, val, initialHashIndex))
	}

	return nil
}

func (c *InMemoryCache[K, V]) checkCacheSize() {
	if c.Size/int(c.capacity) >= int(c.resizeThreshold) {
		c.increaseCacheSize()
	}
}

func (c *InMemoryCache[K, V]) increaseCacheSize() {
	newCache := c.createLargerCache()
	newCache = c.copyValuesToLargerCache(newCache)
	c.setCache(newCache)
}

func (c *InMemoryCache[K, V]) createLargerCache() []*cacheEntry[K, V] {
	return make([]*cacheEntry[K, V], c.resizeCoefficient*c.capacity)
}

func (c *InMemoryCache[K, V]) copyValuesToLargerCache(newCache []*cacheEntry[K, V]) []*cacheEntry[K, V] {
	c.mux.RLock()
	defer c.mux.RUnlock()
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

	return newCache
}

func (c *InMemoryCache[K, V]) setCache(newCache []*cacheEntry[K, V]) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache = newCache
}

func (c *InMemoryCache[K, V]) findNextEmptySpotInCache(key K, index uint32) uint32 {
	var x uint32 = 1
	c.mux.RLock()
	defer c.mux.RUnlock()
	for entry := c.cache[index]; entry != nil && entry.Key != key && !entry.Deleted; entry = c.cache[index] {
		c.checkCapacity(x)

		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	return index
}

func (c *InMemoryCache[K, V]) updateCacheEntry(entry *cacheEntry[K, V], val V) {
	c.mux.Lock()
	defer c.mux.Unlock()
	// update the existing entry
	entry.Val = val
	entry.Deleted = false
}

func (c *InMemoryCache[K, V]) insertNewCacheEntry(index uint32, entry *cacheEntry[K, V]) {
	c.mux.Lock()
	defer c.mux.Unlock()
	// insert the new entry
	c.cache[index] = entry
	c.Size += 1
}

func (c *InMemoryCache[K, V]) createNewCacheEntry(key K, val V, initialHashIndex uint32) *cacheEntry[K, V] {
	return &cacheEntry[K, V]{
		Key:              key,
		Val:              val,
		InitialHashIndex: initialHashIndex,
		Deleted:          false,
	}
}

func (c *InMemoryCache[K, V]) Read(key K) (V, bool) {
	// get the starting index
	index, err := c.hash(key)
	if err != nil {
		panic(err)
	}

	entry := c.findValueInCache(key, index)

	if entry != nil {
		return entry.Val, true
	} else {
		var noop V
		return noop, false
	}
}

func (c *InMemoryCache[K, V]) Remove(key K) error {
	// get the initial index
	index, err := c.hash(key)
	if err != nil {
		return err
	}

	entry := c.findValueInCache(key, index)

	if entry != nil {
		c.deleteEntry(entry)
	}

	return nil
}

func (c *InMemoryCache[K, V]) findValueInCache(key K, index uint32) *cacheEntry[K, V] {
	// find the value with the matching key
	var x uint32 = 1
	c.mux.RLock()
	defer c.mux.RUnlock()

	for entry := c.cache[index]; entry == nil || entry.Key != key; entry = c.cache[index] {
		c.checkCapacity(x)

		// found a nil entry before the key
		if entry == nil {
			return nil
		}

		// find the next index
		index = (index + c.probing(x)) % c.capacity
		x += 1
	}

	return c.cache[index]
}

func (c *InMemoryCache[K, V]) deleteEntry(entry *cacheEntry[K, V]) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry.Deleted = true
	c.Size -= 1
}

func (c *InMemoryCache[K, V]) checkCapacity(x uint32) {
	// there's no space in the cache
	// this can happen if the resizeCoefficient is >= 1
	if x >= c.capacity {
		panic("Cache capacity reached")
	}
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
