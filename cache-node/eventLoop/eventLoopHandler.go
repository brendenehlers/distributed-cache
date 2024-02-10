package eventLoop

import (
	"encoding/json"
	"io"
	"net/http"
)

const INVALID_REQUEST_METHOD string = "Invalid request method"

type RequestBody struct {
	Key string `json:"key"`
	Val []byte `json:"val"`
}

type HandlerWithEventLoop struct {
	El      *EventLoop
	Handler http.Handler
}

func NewHandlerWithEventLoop() *HandlerWithEventLoop {
	el := New()
	el.Run()

	h := http.NewServeMux()
	h.HandleFunc("POST /get", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, el, get)
	})
	h.HandleFunc("POST /set", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, el, set)
	})
	h.HandleFunc("POST /delete", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, el, delete)
	})
	h.HandleFunc("/", invalidResponse)

	return &HandlerWithEventLoop{
		El:      el,
		Handler: h,
	}
}

func handle(w http.ResponseWriter, r *http.Request, el *EventLoop, fn func(el *EventLoop, body RequestBody) ([]byte, error)) {
	body, err := readBody(r.Body)
	if err != nil {
		invalidResponse(w, r)
	}

	val, err := fn(el, body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(val)
	return
}

func get(el *EventLoop, body RequestBody) ([]byte, error) {
	responseChan := make(chan []byte)
	errorChan := make(chan error)

	e := Event{
		Type:         "get",
		Key:          body.Key,
		ResponseChan: responseChan,
		ErrorChan:    errorChan,
	}

	el.Send(e)

	select {
	case val := <-responseChan:
		return val, nil
	case err := <-errorChan:
		return nil, err
	}
}

func set(el *EventLoop, body RequestBody) ([]byte, error) {
	panic("Not implemented")
}

func delete(el *EventLoop, body RequestBody) ([]byte, error) {
	panic("Not implemented")
}

func invalidResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid request"))
}

func validateMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func readBody(reader io.Reader) (RequestBody, error) {
	var body RequestBody
	err := json.NewDecoder(reader).Decode(&body)
	if err != nil {
		return RequestBody{}, err
	}
	return body, nil
}
