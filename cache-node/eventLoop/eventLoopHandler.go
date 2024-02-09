package eventLoop

import "net/http"

type HandlerWithEventLoop struct {
	el *EventLoop
}

func NewHandlerWithEventLoop() *HandlerWithEventLoop {
	el := New()
	el.Run()

	return &HandlerWithEventLoop{
		el: el,
	}
}

func (handler *HandlerWithEventLoop) GetHandler(w http.ResponseWriter, r *http.Request) {
	if ok := handler.validateMethod(r, http.MethodPost); !ok {
		handler.InvalidResponse(w, r)
		return
	}

	w.Write([]byte("Get handler"))
}

func (handler *HandlerWithEventLoop) SetHandler(w http.ResponseWriter, r *http.Request) {
	if ok := handler.validateMethod(r, http.MethodPost); !ok {
		handler.InvalidResponse(w, r)
		return
	}

	w.Write([]byte("Set handler"))
}

func (handler *HandlerWithEventLoop) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if ok := handler.validateMethod(r, http.MethodPost); !ok {
		handler.InvalidResponse(w, r)
		return
	}

	w.Write([]byte("Delete handler"))
}

func (handler *HandlerWithEventLoop) InvalidResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid request"))
}

func (handler *HandlerWithEventLoop) validateMethod(r *http.Request, method string) bool {
	return r.Method == method
}
