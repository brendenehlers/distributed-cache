package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEvent(t *testing.T) {
	expectedType := GET_EVENT_KEY
	expectedKey := "key"

	event, _, _ := CreateGetEvent(expectedKey)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
}

func TestSetEvent(t *testing.T) {
	expectedType := SET_EVENT_KEY
	expectedKey := "key"
	var expectedValue CacheEntry = "value"

	event, _, _ := CreateSetEvent(expectedKey, expectedValue)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
	assert.Equal(t, expectedValue, event.Val)
}

func TestDeleteEvent(t *testing.T) {
	expectedType := DELETE_EVENT_KEY
	expectedKey := "key"

	event, _, _ := CreateDeleteEvent(expectedKey)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
}

func TestSendResponse(t *testing.T) {
	event, respChan, errChan := newEvent(GET_EVENT_KEY, "test", "value")

	event.sendResponse(CacheEventResponse{
		Ok: true,
	})

	select {
	case resp := <-respChan:
		assert.True(t, resp.Ok)
		break
	case err := <-errChan:
		t.Fatalf("error occurred: %s", err)
		break
	default:
		t.Fatal("no value present in responseChan")
	}
}

func TestSendError(t *testing.T) {
	event, respChan, errChan := newEvent(GET_EVENT_KEY, "test", "value")

	event.sendError(fmt.Errorf("test error"))

	select {
	case err := <-errChan:
		if err.Error() != "test error" {
			t.Fatalf("error occurred: %s", err)
		}
		break
	case <-respChan:
		t.Fatalf("response in error test")
		break
	default:
		t.Fatal("no value present in errorChan")
	}
}

func TestCreateEventResponse(t *testing.T) {
	var expectedVal CacheEntry = "test"
	eventResp := createEventResponse(true, expectedVal)

	assert.True(t, eventResp.Ok)
	assert.Equal(t, expectedVal, eventResp.Value)
}
