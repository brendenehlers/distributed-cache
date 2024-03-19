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

const (
	SUCCESS_KEY   = "success"
	SUCCESS_VALUE = "my value"
	ERROR_KEY     = "error"
)

type MockEventLoop struct {
	events chan *loop.CacheEvent
}

func (el *MockEventLoop) Send(event *loop.CacheEvent) {
	el.events <- event

	if event.Key == SUCCESS_KEY {
		<-el.events
		switch event.Type {
		case loop.GET_EVENT_KEY:
			event.ResponseChan <- loop.CacheEventResponse{
				Ok:    true,
				Value: SUCCESS_VALUE,
			}
		case loop.SET_EVENT_KEY:
			fallthrough
		case loop.DELETE_EVENT_KEY:
			event.ResponseChan <- loop.CacheEventResponse{
				Ok: true,
			}
		default:
			panic("invalid event type")
		}

		return
	}

	if event.Key == ERROR_KEY {
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

func TestHandleGet(t *testing.T) {
	key := SUCCESS_KEY
	server := createServerWithEventLoop()

	resp, err := server.handleGetEvent(key)
	if err != nil {
		handleError(t, err)
	}

	assert.NotNil(t, resp)
	assert.True(t, resp.Ok)
	assert.Equal(t, SUCCESS_VALUE, resp.Value)
}

func TestHandleGetError(t *testing.T) {
	key := ERROR_KEY
	server := createServerWithEventLoop()

	_, err := server.handleGetEvent(key)

	assert.NotNil(t, err)
}

func TestHandleSet(t *testing.T) {
	key := SUCCESS_KEY
	value := "test"
	server := createServerWithEventLoop()

	resp, err := server.handleSetEvent(key, value)
	if err != nil {
		handleError(t, err)
	}

	assert.NotNil(t, resp)
	assert.True(t, resp.Ok)
}

func TestHandleSetError(t *testing.T) {
	key := ERROR_KEY
	value := "test"
	server := createServerWithEventLoop()

	_, err := server.handleSetEvent(key, value)

	assert.NotNil(t, err)
}

func TestHandleDelete(t *testing.T) {
	key := SUCCESS_KEY
	server := createServerWithEventLoop()

	resp, err := server.handleDeleteEvent(key)
	if err != nil {
		handleError(t, err)
	}

	assert.NotNil(t, resp)
	assert.True(t, resp.Ok)
}

func TestHandleDeleteError(t *testing.T) {
	key := ERROR_KEY
	server := createServerWithEventLoop()

	_, err := server.handleDeleteEvent(key)

	assert.NotNil(t, err)
}

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

func TestSendEvent(t *testing.T) {
	server := createServerWithEventLoop()
	event, r, e := loop.CreateGetEvent(SUCCESS_KEY)

	resp, err := server.sendEvent(event, r, e)
	if err != nil {
		handleError(t, err)
	}

	assert.NotNil(t, resp)
	assert.True(t, resp.Ok)
	assert.Equal(t, SUCCESS_VALUE, resp.Value)
}

func TestSendEventError(t *testing.T) {
	el := createMockEventLoop()
	server := createServer(el)
	event, r, e := loop.CreateGetEvent(ERROR_KEY)
	_, err := server.sendEvent(event, r, e)

	assert.NotNil(t, err)
}

func TestEncodeResponse(t *testing.T) {
	expectedMsg := "test response"
	resp := Response{
		Message: expectedMsg,
	}

	enc, err := encodeResponse(resp)

	assert.Nil(t, err)
	assert.Contains(t, enc.String(), expectedMsg)
}

func TestCreateErrorResponse(t *testing.T) {
	errMsg := "my fancy error"
	err := fmt.Errorf(errMsg)
	resp := createErrorResponse(err)

	assert.Equal(t, ERROR_MSG, resp.Message)
	assert.Equal(t, errMsg, resp.Error)
}

func TestCreateGetResponseOk(t *testing.T) {
	var val loop.CacheEntry = "test"
	resp := createGetResponse(true, val)

	assert.NotNil(t, resp.Value)
	assert.Equal(t, VALUE_FOUND_MSG, resp.Message)
	assert.Equal(t, val, resp.Value)
}

func TestCreateGetResponseNotOk(t *testing.T) {
	var val loop.CacheEntry = "test"
	resp := createGetResponse(false, val)

	assert.Nil(t, resp.Value)
	assert.Equal(t, VALUE_NOT_FOUND_MSG, resp.Message)
}

func createServer(el EventLoop) *Server {
	return &Server{
		eventLoop: el,
	}
}

func createServerWithEventLoop() *Server {
	el := createMockEventLoop()
	return &Server{
		eventLoop: el,
	}
}

func createMockEventLoop() *MockEventLoop {
	return &MockEventLoop{
		events: make(chan *loop.CacheEvent, 1),
	}
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
