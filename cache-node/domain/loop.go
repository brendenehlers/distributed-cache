package domain

func NewEventLoop(cache Cache) *EventLoopImpl {
	return &EventLoopImpl{
		cache: cache,
		send:  make(chan Event),
		quit:  make(chan int),
	}
}

func (eventLoop *EventLoopImpl) Run() {
	go eventLoop.startLoop()
}

func (eventLoop *EventLoopImpl) startLoop() {
	for {
		select {
		case event := <-eventLoop.send:
			eventLoop.handleEvent(event)
		case <-eventLoop.quit:
			return
		}
	}
}

func (eventLoop *EventLoopImpl) handleEvent(event Event) {
	switch event.Type {
	case "get":
		go eventLoop.getValue(event)
	case "set":
		go eventLoop.setValue(event)
	case "delete":
		go eventLoop.deleteValue(event)
	default:
		panic("Invalid event type")
	}
}

func (eventLoop *EventLoopImpl) getValue(event Event) {
	val, ok := eventLoop.cache.Get(string(event.Key))
	writeResponseToChannel(event.ResponseChan, ok, val)
}

func (eventLoop *EventLoopImpl) setValue(event Event) {
	if err := eventLoop.cache.Set(event.Key, event.Val); checkError(err) {
		writeErrorToChannel(event.ErrorChan, err)
	} else {
		writeResponseToChannel(event.ResponseChan, true, nil)
	}
}

func (eventLoop *EventLoopImpl) deleteValue(event Event) {
	if err := eventLoop.cache.Delete(event.Key); checkError(err) {
		writeErrorToChannel(event.ErrorChan, err)
	} else {
		writeResponseToChannel(event.ResponseChan, true, nil)
	}
}

func checkError(err error) bool {
	return err != nil
}

func writeErrorToChannel(channel chan error, err error) {
	channel <- err
}

func writeResponseToChannel(channel chan EventResponse, ok bool, value CacheEntry) {
	channel <- EventResponse{
		Ok:    ok,
		Value: value,
	}
}

func (eventLoop *EventLoopImpl) Stop() {
	eventLoop.quit <- 1
}

func (eventLoop *EventLoopImpl) Send(e Event) {
	eventLoop.send <- e
}
