package domain

func NewEventLoop(cache Cache) *EventLoop {
	return &EventLoop{
		cache: cache,
		send:  make(chan Event),
		quit:  make(chan int),
	}
}

func (eventLoop *EventLoop) Run() {
	go eventLoop.startLoop()
}

func (eventLoop *EventLoop) startLoop() {
	for {
		select {
		case event := <-eventLoop.send:
			eventLoop.handleEvent(event)
		case <-eventLoop.quit:
			return
		}
	}
}

func (eventLoop *EventLoop) handleEvent(event Event) {
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

func (eventLoop *EventLoop) getValue(event Event) {
	val, ok := eventLoop.cache.Get(string(event.Key))
	writeResponseToChannel(event.ResponseChan, ok, val)
}

func (eventLoop *EventLoop) setValue(event Event) {
	if err := eventLoop.cache.Set(event.Key, event.Val); checkError(err) {
		writeErrorToChannel(event.ErrorChan, err)
	} else {
		writeResponseToChannel(event.ResponseChan, true, nil)
	}
}

func (eventLoop *EventLoop) deleteValue(event Event) {
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

func (eventLoop *EventLoop) Stop() {
	eventLoop.quit <- 1
}

func (eventLoop *EventLoop) Send(e Event) {
	eventLoop.send <- e
}
