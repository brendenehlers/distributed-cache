package domain

import "fmt"

func NewEventLoop(cache Cache) *EventLoop {
	return &EventLoop{
		cache: cache,
		send:  make(chan Event),
		quit:  make(chan int),
	}
}

func (eventLoop *EventLoop) Run() {
	go func() {
		for {
			select {
			case event := <-eventLoop.send:
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
			case <-eventLoop.quit:
				return
			}
		}
	}()
}

func (eventLoop *EventLoop) getValue(event Event) {
	fmt.Printf("Getting %s from the cache\n", event.Key)

	val, ok := eventLoop.cache.Get(string(event.Key))

	writeResponseToChannel(event.ResponseChan, ok, val)
}

func (eventLoop *EventLoop) setValue(event Event) {
	fmt.Printf("Setting %s to %s in the cache\n", event.Key, event.Val)

	err := eventLoop.cache.Set(event.Key, event.Val)
	if err != nil {
		event.ErrorChan <- err
		return
	}

	writeResponseToChannel(event.ResponseChan, true, nil)
}

func (eventLoop *EventLoop) deleteValue(event Event) {
	fmt.Printf("Deleting %s from the cache\n", event.Key)

	err := eventLoop.cache.Delete(event.Key)
	if err != nil {
		event.ErrorChan <- err
		return
	}

	writeResponseToChannel(event.ResponseChan, true, nil)
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
