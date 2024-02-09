package eventLoop

type EventLoop struct {
	quit chan struct{}
}

func New() *EventLoop {
	return &EventLoop{
		quit: make(chan struct{}),
	}
}
