package domain

func CreateGetEvent(key string) (event Event, responseChan chan EventResponse, errorChan chan error) {
	return newEvent("get", key, nil)
}

func CreateSetEvent(key string, value CacheEntry) (event Event, responseChan chan EventResponse, errorChan chan error) {
	return newEvent("set", key, value)
}

func CreateDeleteEvent(key string) (event Event, responseChan chan EventResponse, errorChan chan error) {
	return newEvent("delete", key, nil)
}

func newEvent(eventType string, key string, value CacheEntry) (event Event, responseChan chan EventResponse, errorChan chan error) {
	responseChan = make(chan EventResponse)
	errorChan = make(chan error)
	event = Event{
		Type:         eventType,
		Key:          key,
		Val:          value,
		ResponseChan: responseChan,
		ErrorChan:    errorChan,
	}

	return event, responseChan, errorChan
}
