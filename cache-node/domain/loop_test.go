package domain

import (
	"testing"
)

func TestGetValue(t *testing.T) {
	eventLoop := setup()
	key := "test"
	value := "value"
	eventLoop.cache.Set(key, value)

	event, responseChan, errorChan := CreateGetEvent(key)

	go eventLoop.getValue(event)

	var resp EventResponse
	select {
	case resp = <-responseChan:
		if !resp.Ok {
			t.Fatalf("response value from cache was not okay")
		}
		if resp.Value != value {
			t.Fatalf("expected value (%s) != actual value (%s)", value, resp.Value)
		}
		break
	case err := <-errorChan:
		t.Fatalf("test failed with error: %s", err)
		break
	}

	if !resp.Ok || resp.Value == "" {
		t.Fatalf("resp object was not defined")
	}
}

func TestGetNoValue(t *testing.T) {
	eventLoop := setup()
	key := "test"

	event, responseChan, errorChan := CreateGetEvent(key)

	go eventLoop.getValue(event)

	var resp EventResponse
	select {
	case resp = <-responseChan:
		if resp.Ok {
			t.Fatalf("response value from cache was okay when no value set")
		}
		break
	case err := <-errorChan:
		t.Fatalf("test failed with error: %s", err)
		break
	}
}

func TestSetValue(t *testing.T) {
	eventLoop := setup()
	key := "key"
	value := "value"
	eventLoop.cache.Set(key, value)

	event, responseChan, errorChan := CreateSetEvent(key, value)

	go eventLoop.setValue(event)

	select {
	case resp := <-responseChan:
		if !resp.Ok {
			t.Fatal("response was not okay")
		}
		break
	case err := <-errorChan:
		t.Fatalf("test failed with error: %s", err)
		break
	}

	val, ok := eventLoop.cache.Get(key)
	if !ok {
		t.Fatal("Getting from cache was not okay")
	}

	if val != value {
		t.Fatalf("expected (%s) != actual (%s)", value, val)
	}
}

func TestDeleteValue(t *testing.T) {
	eventLoop := setup()
	key := "key"
	value := "value"
	eventLoop.cache.Set(key, value)

	event, responseChan, errorChan := CreateDeleteEvent(key)

	go eventLoop.deleteValue(event)

	select {
	case resp := <-responseChan:
		if !resp.Ok {
			t.Fatal("response was not okay")
		}
		break
	case err := <-errorChan:
		t.Fatalf("test failed with error: %s", err)
		break
	}

	_, ok := eventLoop.cache.Get(key)
	if ok {
		t.Fatal("Getting from cache was okay when shouldn't be")
	}
}

type MockCache struct {
	cache map[string]CacheEntry
}

func (mc *MockCache) Get(key string) (CacheEntry, bool) {
	val, ok := mc.cache[key]
	return val, ok
}

func (mc *MockCache) Set(key string, val CacheEntry) error {
	mc.cache[key] = val
	return nil
}

func (mc *MockCache) Delete(key string) error {
	delete(mc.cache, key)
	return nil
}

func setup() *EventLoopImpl {
	mockCache := &MockCache{
		cache: make(map[string]CacheEntry),
	}
	eventLoop := NewEventLoop(mockCache)

	return eventLoop
}
