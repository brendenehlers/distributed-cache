package domain

type CacheEntry any

type Event struct {
	Type         string
	Key          string
	Val          CacheEntry
	ResponseChan chan EventResponse
	ErrorChan    chan error
}

type EventResponse struct {
	Ok    bool
	Value CacheEntry
}

type EventLoopImpl struct {
	cache Cache
	send  chan Event
	quit  chan int
}
