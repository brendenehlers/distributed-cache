package eventLoop

import (
	"fmt"

	"github.com/brendenehlers/go-distributed-cache/cache-node/cache"
)

const SEND_CHANNEL_SIZE = 100

type Event struct {
	Type         string
	Key          string
	Val          []byte
	ResponseChan chan []byte
	ErrorChan    chan error
}

type EventLoop struct {
	cache cache.Cache[string, []byte]
	send  chan Event
	quit  chan int
}

func New() *EventLoop {
	return &EventLoop{
		cache: cache.New[string, []byte](cache.Options{}),
		send:  make(chan Event, SEND_CHANNEL_SIZE),
		quit:  make(chan int),
	}
}

func (el *EventLoop) Run() {

	go func() {
		for {
			select {
			case e := <-el.send:
				// handle event
				switch e.Type {
				case "get":
					fmt.Printf("Getting %s from the cache\n", e.Key)

					val, err := el.cache.Read(string(e.Key))
					if err != nil {
						e.ErrorChan <- err
					}

					e.ResponseChan <- val
				case "set":
				case "delete":
				default:
					panic("Invalid event type")
				}
			case <-el.quit:
				return
			}
		}
	}()

}

func (el *EventLoop) Stop() {
	el.quit <- 1
}

func (el *EventLoop) Send(e Event) {
	el.send <- e
}
