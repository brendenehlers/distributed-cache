package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEvent(t *testing.T) {
	expectedType := "get"
	expectedKey := "key"

	event, _, _ := CreateGetEvent(expectedKey)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
}

func TestSetEvent(t *testing.T) {
	expectedType := "set"
	expectedKey := "key"
	var expectedValue CacheEntry = "value"

	event, _, _ := CreateSetEvent(expectedKey, expectedValue)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
	assert.Equal(t, expectedValue, event.Val)
}

func TestDeleteEvent(t *testing.T) {
	expectedType := "delete"
	expectedKey := "key"

	event, _, _ := CreateDeleteEvent(expectedKey)

	assert.Equal(t, expectedType, event.Type)
	assert.Equal(t, expectedKey, event.Key)
}
