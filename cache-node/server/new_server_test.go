package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
	"github.com/stretchr/testify/assert"
)

type MockEventLoop struct {
	events chan *loop.CacheEvent
}

func (el *MockEventLoop) Send(event *loop.CacheEvent) {
	el.events <- event

	if event.Key == "success" {
		<-el.events
		event.ResponseChan <- loop.CacheEventResponse{
			Ok:    true,
			Value: "my value",
		}
		return
	}

	if event.Key == "error" {
		<-el.events
		event.ErrorChan <- fmt.Errorf("something went wrong")
	}
}

func (el *MockEventLoop) Run() {}

func (el *MockEventLoop) Stop() {}

/**
Test Cases:
createMux returns the correct handler
newServer creates a new server object and returns it
server handles get/set/delete requests
get/set/delete use POST method

handleGet calls the cache to get the value
get returns the correct value from the cache
the handleSet method takes a http.Request and returns a Response
set inserts/updates the correct value in the cache
the handleDelete method takes a http.Request and returns a Response
delete removes the correct value in the cache
*/

// func TestHandleGetReturnsCacheEventResponse(t *testing.T) {
// 	body, _ := createReqBody("test", nil)
// 	resp, err := handleGet(&MockEventLoop{}, &http.Request{Body: body})
// 	if err != nil {
// 		handleError(t, err)
// 	}
// 	assert.IsType(t, loop.CacheEventResponse{}, resp)
// }

// func TestGetHandlerSendsEventToLoop(t *testing.T) {
// 	eventLoop := &MockEventLoop{events: make(chan *loop.CacheEvent, 1)}

// 	expectedKey := "test"
// 	reqBody := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"key": "%s"}`, expectedKey)))
// 	_, err := handleGet(eventLoop, &http.Request{Body: reqBody})
// 	if err != nil {
// 		handleError(t, err)
// 	}

// 	assert.True(t, len(eventLoop.events) == 1)

// 	select {
// 	case event := <-eventLoop.events:
// 		assert.Equal(t, expectedKey, event.Key)
// 		break
// 	default:
// 		t.Fatal("No value present in events channel")
// 	}
// }

func TestReadRequestBody(t *testing.T) {
	expectedKey, expectedValue := "test", "my value"
	r, err := createReqBody(expectedKey, expectedValue)
	if err != nil {
		handleError(t, err)
	}

	data := RequestBody{}
	if err = decodeRequestBody(r, &data); err != nil {
		handleError(t, err)
	}

	assert.Equal(t, expectedKey, data.Key)
	assert.Equal(t, expectedValue, data.Value)
}

func createReqBody(key string, value any) (io.ReadCloser, error) {
	data := RequestBody{
		Key:   key,
		Value: value,
	}
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(&data)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(buf), nil
}

func handleError(t *testing.T, err error) {
	t.Fatalf("error occurred: %s", err)
	t.FailNow()
}
