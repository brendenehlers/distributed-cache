package cache

import (
	"testing"
)

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
		val, err := cache.Read(k)
		if err != nil {
			t.Fatalf("TestRead: failed on read with err: %s\n", err)
		}

		if mapVal := kvs[k]; val != mapVal {
			t.Fatalf("TestRead: cache val (%s) != map val (%s)\n", val, mapVal)
		}
	}
}
