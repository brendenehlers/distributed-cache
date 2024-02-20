package cache

import (
	"math/rand"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	cache := New[int, int](Options{})

	kvs := make(map[int]int)

	keys := make([]int, 10000)
	for _, key := range keys {
		kvs[key] = rand.Intn(10000)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for k, v := range kvs {
			cache.Insert(k, v)
		}
	}
}

func TestInsert(t *testing.T) {

	cache := New[int, string](Options{})

	kvs := make(map[int]string)
	kvs[123] = "hello world 1"
	kvs[234] = "hello world 2"
	kvs[345] = "hello world 3"
	kvs[456] = "hello world 4"

	for k, v := range kvs {
		if err := cache.Insert(k, v); err != nil {
			t.Fatalf("TestInsert: failed on insert with err: %s\n", err)
		}
	}

	if len(kvs) != cache.Size {
		t.Fatalf("TestInsert: Initial elements length (%d) != cache size (%d)\n", len(kvs), cache.Size)
	}
}

func TestRead(t *testing.T) {
	cache := New[int, string](Options{})

	kvs := make(map[int]string)
	kvs[123] = "hello world 1"
	kvs[234] = "hello world 2"
	kvs[345] = "hello world 3"
	kvs[456] = "hello world 4"

	for k, v := range kvs {
		if err := cache.Insert(k, v); err != nil {
			t.Fatalf("TestRead: failed on insert with err: %s\n", err)
		}
	}

	for k := range kvs {
		val, ok := cache.Read(k)
		if !ok {
			t.Fatalf("TestRead: failed on read ")
		}

		if mapVal := kvs[k]; val != mapVal {
			t.Fatalf("TestRead: cache val (%s) != map val (%s)\n", val, mapVal)
		}
	}
}

func TestRemove(t *testing.T) {
	cache := New[int, string](Options{})

	kvs := make(map[int]string)
	kvs[123] = "hello world 1"
	kvs[234] = "hello world 2"
	kvs[345] = "hello world 3"
	kvs[456] = "hello world 4"

	for k, v := range kvs {
		if err := cache.Insert(k, v); err != nil {
			t.Fatalf("TestRemove: failed on insert with err: %s\n", err)
		}
	}

	for k := range kvs {
		if err := cache.Remove(k); err != nil {
			t.Fatalf("TestRemove: failed on remove with err: %s\n", err)
		}
	}

	if cache.Size != 0 {
		t.Fatal("Cache has non-zero size after removing all values")
	}

}
