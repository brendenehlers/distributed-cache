package eventLoop

type EventLoop struct {
	quit chan int
}

func New() *EventLoop {
	return &EventLoop{
		quit: make(chan int),
	}
}

func (el *EventLoop) Run() {
	close(el.quit)
}
