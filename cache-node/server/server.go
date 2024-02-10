package server

import (
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/eventLoop"
)

func NewServer() *http.Server {

	elh := eventLoop.NewHandlerWithEventLoop()

	// log incoming requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s: %s %s\n", r.RemoteAddr, r.Method, r.URL)
		log.Printf("User Agent: %s\n", r.UserAgent())

		elh.Handler.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return s
}
