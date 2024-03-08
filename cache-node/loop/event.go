package loop

const (
	GET_EVENT_KEY    = "get"
	SET_EVENT_KEY    = "set"
	DELETE_EVENT_KEY = "delete"
)

type CacheEvent struct {
	Type         string
	Key          string
	Val          CacheEntry
	responseChan chan CacheEventResponse
	errorChan    chan error
}

type CacheEventResponse struct {
	Ok    bool
	Value CacheEntry
}

func CreateGetEvent(key string) (event *CacheEvent, responseChan chan CacheEventResponse, errorChan chan error) {
	return newEvent(GET_EVENT_KEY, key, nil)
}

func CreateSetEvent(key string, value CacheEntry) (event *CacheEvent, responseChan chan CacheEventResponse, errorChan chan error) {
	return newEvent(SET_EVENT_KEY, key, value)
}

func CreateDeleteEvent(key string) (event *CacheEvent, responseChan chan CacheEventResponse, errorChan chan error) {
	return newEvent(DELETE_EVENT_KEY, key, nil)
}

func newEvent(eventType string, key string, value CacheEntry) (event *CacheEvent, responseChan chan CacheEventResponse, errorChan chan error) {
	responseChan = make(chan CacheEventResponse, 1)
	errorChan = make(chan error, 1)
	event = &CacheEvent{
		Type:         eventType,
		Key:          key,
		Val:          value,
		responseChan: responseChan,
		errorChan:    errorChan,
	}

	return event, responseChan, errorChan
}

func createEventResponse(ok bool, value CacheEntry) CacheEventResponse {
	return CacheEventResponse{
		Ok:    ok,
		Value: value,
	}
}

func (event *CacheEvent) sendResponse(resp CacheEventResponse) {
	event.responseChan <- resp
}

func (event *CacheEvent) sendError(err error) {
	event.errorChan <- err
}
