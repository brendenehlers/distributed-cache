package domain

import "testing"

func TestGetEvent(t *testing.T) {
	eventType := "get"
	expected_key := "key"

	event, _, _ := GetEvent(expected_key)

	assertActualEqualExpected(t, event.Type, eventType)
	assertActualEqualExpected(t, event.Key, expected_key)
	assertActualEqualExpected(t, event.Val, nil)
}

func TestSetEvent(t *testing.T) {
	eventType := "set"
	expected_key := "key"
	var expected_value CacheEntry = "value"

	event, _, _ := SetEvent(expected_key, expected_value)

	assertActualEqualExpected(t, event.Type, eventType)
	assertActualEqualExpected(t, event.Key, expected_key)
	assertActualEqualExpected(t, event.Val, expected_value)
}

func TestDeleteEvent(t *testing.T) {
	eventType := "delete"
	expected_key := "key"

	event, _, _ := DeleteEvent(expected_key)

	assertActualEqualExpected(t, event.Type, eventType)
	assertActualEqualExpected(t, event.Key, expected_key)
	assertActualEqualExpected(t, event.Val, nil)
}

func assertActualEqualExpected[T comparable](t *testing.T, actual T, expected T) {
	if actual != expected {
		t.Fatalf("actual (%v) != expected (%v)", actual, expected)
	}
}
