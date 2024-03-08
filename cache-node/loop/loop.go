package loop

const (
	DEFAULT_EVENTS_CHANNEL_CAP = 50
	DEFAULT_QUIT_CHANNEL_CAP   = 1

	PROCESSED_EVENT_CODE = 0
	KILL_CODE            = 1
)

type EventLoopImpl struct {
	cache  Cache
	events chan *CacheEvent
	quit   chan int
}

func NewEventLoop(cache Cache) *EventLoopImpl {
	return &EventLoopImpl{
		cache:  cache,
		events: make(chan *CacheEvent, DEFAULT_EVENTS_CHANNEL_CAP),
		quit:   make(chan int, DEFAULT_QUIT_CHANNEL_CAP),
	}
}

func (eventLoop *EventLoopImpl) Send(event *CacheEvent) {
	eventLoop.events <- event
}

func (eventLoop *EventLoopImpl) Stop() {
	eventLoop.quit <- 1
}

func (eventLoop *EventLoopImpl) Run() {
	// this function is untested but the component methods are well tested
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
	event.sendResponse(createEventResponse(ok, value))
}

func (eventLoop *EventLoopImpl) handleSetEvent(event *CacheEvent) {
	err := eventLoop.cache.Set(event.Key, event.Val)
	if err != nil {
		event.sendError(err)
		return
	}
	event.sendResponse(createEventResponse(true, nil))
}

func (eventLoop *EventLoopImpl) handleDeleteEvent(event *CacheEvent) {
	err := eventLoop.cache.Delete(event.Key)
	if err != nil {
		event.sendError(err)
		return
	}
	event.sendResponse(createEventResponse(true, nil))
}
