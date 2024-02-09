package server

import (
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/eventLoop"
)

type EventLoopHandler struct {
	el *eventLoop.EventLoop
}

func NewServer() *http.Server {
	h := http.NewServeMux()

	elh := &EventLoopHandler{
		el: eventLoop.New(),
	}

	h.HandleFunc("/get", elh.GetHandler)
	h.HandleFunc("/set", elh.SetHandler)
	h.HandleFunc("/delete", elh.DeleteHandler)
	h.HandleFunc("/", invalidResponse)

	// log incoming requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s: %s %s\n", r.RemoteAddr, r.Method, r.URL)
		log.Printf("User Agent: %s\n", r.UserAgent())

		h.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return s
}

func validateMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func (el *EventLoopHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Get handler"))
}

func (el *EventLoopHandler) SetHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Set handler"))
}

func (el *EventLoopHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Delete handler"))
}

func invalidResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid request"))
}
