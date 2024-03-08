package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCache struct {
	cache map[string]CacheEntry
}

func (mc *MockCache) Get(key string) (CacheEntry, bool) {
	if key == "miss" {
		return nil, false
	}

	val, ok := mc.cache[key]
	return val, ok
}

func (mc *MockCache) Set(key string, val CacheEntry) error {
	if key == "error" {
		return fmt.Errorf("error setting value")
	}
	mc.cache[key] = val
	return nil
}

func (mc *MockCache) Delete(key string) error {
	if key == "error" {
		return fmt.Errorf("error deleting value")
	}
	mc.cache[key] = nil
	return nil
}

func TestNewEventLoopMethodReturnsEventLoop(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	assert.NotNil(t, eventLoop)
}

func TestNewEventLoopHasDefaults(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	assert.NotNil(t, eventLoop.cache)
	assert.NotNil(t, eventLoop.events)
	assert.NotNil(t, eventLoop.quit)
}

func TestDefaultEventsChannelCap(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	assert.Equal(t, DEFAULT_EVENTS_CHANNEL_CAP, cap(eventLoop.events))
}

func TestDefaultQuitChannelCap(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	assert.Equal(t, DEFAULT_QUIT_CHANNEL_CAP, cap(eventLoop.quit))
}

func TestStopMethodSendsMessageOnQuitChannel(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	eventLoop.Stop()

	assert.True(t, len(eventLoop.quit) == 1)
	assert.NotNil(t, <-eventLoop.quit)
}

func TestSendingValueToQuitChannelReturns(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	eventLoop.quit <- 1
	code := eventLoop.multiplexChannels()

	assert.Equal(t, KILL_CODE, code)
}

func TestIsKillCodeReturnsTrueOnKillCode(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	isKillCode := eventLoop.isKillCode(KILL_CODE)
	assert.True(t, isKillCode)
}

func TestIsKillCodeReturnsFalseOnNonKillCode(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	isKillCode := eventLoop.isKillCode(PROCESSED_EVENT_CODE)
	assert.False(t, isKillCode)
}

func TestSendMethodPushesEventToEventsChannel(t *testing.T) {
	eventLoop := createEmptyEventLoop()
	expectedType := "test event"
	eventLoop.Send(&CacheEvent{Type: expectedType})

	assert.True(t, len(eventLoop.events) == 1)

	event := <-eventLoop.events
	assert.NotNil(t, event)
	assert.Equal(t, expectedType, event.Type)
}

func TestSendingEventCallsHandleEvent(t *testing.T) {
	key := "test"
	var val CacheEntry = "my value"
	event, _, _ := CreateSetEvent(key, val)
	eventLoop := createEmptyEventLoop()

	eventLoop.Send(event)
	code := eventLoop.multiplexChannels()

	assert.Equal(t, PROCESSED_EVENT_CODE, code)
}

func TestHandleGetEventSendsValueInResponseChan(t *testing.T) {
	key := "test"
	var expectedValue CacheEntry = "my value"
	event, responseChan, errorChan := CreateGetEvent(key)
	eventLoop := createEmptyEventLoop()
	setCacheValue(eventLoop, key, expectedValue)

	eventLoop.handleGetEvent(event)

	select {
	case resp := <-responseChan:
		assert.True(t, resp.Ok)
		assert.Equal(t, expectedValue, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestHandleGetEventSendsNonOkayResponseOnMiss(t *testing.T) {
	key := "miss"
	event, responseChan, errorChan := CreateGetEvent(key)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleGetEvent(event)

	select {
	case resp := <-responseChan:
		assert.False(t, resp.Ok)
		assert.Nil(t, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestHandleSetEventUpdatesCache(t *testing.T) {
	expectedKey := "test"
	var expectedValue CacheEntry = "my value"
	event, _, _ := CreateSetEvent(expectedKey, expectedValue)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleSetEvent(event)

	err := eventLoop.cache.Set(expectedKey, expectedValue)
	cacheValue := getCacheValue(eventLoop, expectedKey)
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, cacheValue)
}

func TestHandleSetEventSendsSuccessResponse(t *testing.T) {
	key := "test"
	var val CacheEntry = "val"
	event, responseChan, errorChan := CreateSetEvent(key, val)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleSetEvent(event)

	select {
	case resp := <-responseChan:
		assert.NotNil(t, resp)
		assert.True(t, resp.Ok)
		assert.Nil(t, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestHandleSetEventSendErrorResponse(t *testing.T) {
	key := "error"
	var val CacheEntry = "val"
	event, responseChan, errorChan := CreateSetEvent(key, val)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleSetEvent(event)

	select {
	case err := <-errorChan:
		assert.NotNil(t, err)
		break
	case <-responseChan:
		t.Fatal("non-error response during error test")
		break
	default:
		handleDefault(t, "errorChan")
	}
}

func TestHandleDeleteEventDeletesValueFromCache(t *testing.T) {
	key := "test"
	var value CacheEntry = "my value"
	event, _, _ := CreateDeleteEvent(key)
	eventLoop := createEmptyEventLoop()
	setCacheValue(eventLoop, key, value)

	eventLoop.handleDeleteEvent(event)

	assert.Nil(t, getCacheValue(eventLoop, key))
}

func TestHandleDeleteEventSendResponse(t *testing.T) {
	key := "test"
	var value CacheEntry = "my value"
	event, responseChan, errorChan := CreateDeleteEvent(key)
	eventLoop := createEmptyEventLoop()
	setCacheValue(eventLoop, key, value)

	eventLoop.handleDeleteEvent(event)

	select {
	case resp := <-responseChan:
		assert.NotNil(t, resp)
		assert.True(t, resp.Ok)
		assert.Nil(t, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestHandleDeleteEventSendsError(t *testing.T) {
	key := "error"
	event, responseChan, errorChan := CreateDeleteEvent(key)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleDeleteEvent(event)

	select {
	case err := <-errorChan:
		assert.NotNil(t, err)
		break
	case <-responseChan:
		t.Fatal("non-error response during error test")
		break
	default:
		handleDefault(t, "errorChan")
	}
}

func TestCallingHandleEventWithGetEventCallsHandleGetEvent(t *testing.T) {
	key := "test"
	var expectedValue CacheEntry = "my value"
	event, responseChan, errorChan := CreateGetEvent(key)
	eventLoop := createEmptyEventLoop()
	setCacheValue(eventLoop, key, expectedValue)

	eventLoop.handleEvent(event)

	select {
	case resp := <-responseChan:
		assert.True(t, resp.Ok)
		assert.Equal(t, expectedValue, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestCallingHandleEventWithSetEventCallsHandleSetEvent(t *testing.T) {
	key := "test"
	var val CacheEntry = "val"
	event, responseChan, errorChan := CreateSetEvent(key, val)
	eventLoop := createEmptyEventLoop()

	eventLoop.handleEvent(event)

	select {
	case resp := <-responseChan:
		assert.NotNil(t, resp)
		assert.True(t, resp.Ok)
		assert.Nil(t, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestCallingHandleEventWithDeleteEventCallsHandleDeleteEvent(t *testing.T) {
	key := "test"
	var value CacheEntry = "my value"
	event, responseChan, errorChan := CreateDeleteEvent(key)
	eventLoop := createEmptyEventLoop()
	setCacheValue(eventLoop, key, value)

	eventLoop.handleEvent(event)

	select {
	case resp := <-responseChan:
		assert.NotNil(t, resp)
		assert.True(t, resp.Ok)
		assert.Nil(t, resp.Value)
		break
	case err := <-errorChan:
		handleError(t, err)
		break
	default:
		handleDefault(t, "responseChan")
	}
}

func TestCallingHandleEventWithRandomEventPanics(t *testing.T) {
	eventLoop := createEmptyEventLoop()

	defer func(t *testing.T) {
		if r := recover(); r == nil {
			t.Fatal("Panic was nil")
		}
	}(t)
	eventLoop.handleEvent(&CacheEvent{Type: "not defined"})

	t.Fatal("Method call didn't panic")
}

func createEmptyEventLoop() *EventLoopImpl {
	eventLoop := NewEventLoop(&MockCache{
		cache: make(map[string]CacheEntry),
	})
	return eventLoop
}

func setCacheValue(eventLoop *EventLoopImpl, key string, value CacheEntry) {
	eventLoop.cache.(*MockCache).cache[key] = value
}

func getCacheValue(eventLoop *EventLoopImpl, key string) CacheEntry {
	return eventLoop.cache.(*MockCache).cache[key]
}

func handleError(t *testing.T, err error) {
	t.Fatalf("an error occurred: %s", err)
}

func handleDefault(t *testing.T, expectedChan string) {
	t.Fatalf("no value present in %s", expectedChan)
}
