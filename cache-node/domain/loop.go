package domain

const (
	DEFAULT_EVENTS_CHANNEL_CAP = 50
	DEFAULT_QUIT_CHANNEL_CAP   = 1

	PROCESSED_EVENT_CODE = 0
	KILL_CODE            = 1
)

func NewEventLoop(cache Cache) *EventLoopImpl {
	return &EventLoopImpl{
		cache:  cache,
		events: make(chan *CacheEvent, DEFAULT_EVENTS_CHANNEL_CAP),
		quit:   make(chan int, DEFAULT_QUIT_CHANNEL_CAP),
	}
}

type EventLoopImpl struct {
	cache  Cache
	events chan *CacheEvent
	quit   chan int
}

func (eventLoop *EventLoopImpl) Send(event *CacheEvent) {
	eventLoop.events <- event
}

func (eventLoop *EventLoopImpl) Stop() {
	eventLoop.quit <- 1
}

func (eventLoop *EventLoopImpl) Run() {
	// this function is untested
	// but the component methods are well tested, so I figure it's alright
	for {
		code := eventLoop.multiplexChannels()

		if eventLoop.isKillCode(code) {
			return
		}
	}
}

func (eventLoop *EventLoopImpl) multiplexChannels() int {
	select {
	case event := <-eventLoop.events:
		eventLoop.handleEvent(event)
		return PROCESSED_EVENT_CODE
	case <-eventLoop.quit:
		return KILL_CODE
	}
}

func (eventLoop *EventLoopImpl) isKillCode(code int) bool {
	return code == KILL_CODE
}

func (eventLoop *EventLoopImpl) handleEvent(event *CacheEvent) {
	switch event.Type {
	case GET_EVENT_KEY:
		eventLoop.handleGetEvent(event)
	case SET_EVENT_KEY:
		eventLoop.handleSetEvent(event)
	case DELETE_EVENT_KEY:
		eventLoop.handleDeleteEvent(event)
	default:
		panic("unknown event type")
	}
}

func (eventLoop *EventLoopImpl) handleGetEvent(event *CacheEvent) {
	value, ok := eventLoop.cache.Get(event.Key)
	sendResponseOnResponseChan(event.ResponseChan, ok, value)
}

func (eventLoop *EventLoopImpl) handleSetEvent(event *CacheEvent) {
	err := eventLoop.cache.Set(event.Key, event.Val)
	if err != nil {
		sendError(event.ErrorChan, err)
		return
	}

	sendResponseOnResponseChan(event.ResponseChan, true, nil)
}

func (eventLoop *EventLoopImpl) handleDeleteEvent(event *CacheEvent) {
	err := eventLoop.cache.Delete(event.Key)
	if err != nil {
		sendError(event.ErrorChan, err)
		return
	}

	sendResponseOnResponseChan(event.ResponseChan, true, nil)
}

func sendError(channel chan error, err error) {
	channel <- err
}

func sendResponseOnResponseChan(channel chan CacheEventResponse, ok bool, value CacheEntry) {
	channel <- createEventResponse(ok, value)
}

func createEventResponse(ok bool, value CacheEntry) CacheEventResponse {
	return CacheEventResponse{
		Ok:    ok,
		Value: value,
	}
}
