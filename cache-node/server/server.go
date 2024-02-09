package server

import (
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/eventLoop"
)

func NewServer() *http.Server {

	elh := eventLoop.NewHandlerWithEventLoop()

	mux := http.NewServeMux()
	mux.HandleFunc("/get", elh.GetHandler)
	mux.HandleFunc("/set", elh.SetHandler)
	mux.HandleFunc("/delete", elh.DeleteHandler)
	mux.HandleFunc("/", elh.InvalidResponse)

	// log incoming requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s: %s %s\n", r.RemoteAddr, r.Method, r.URL)
		log.Printf("User Agent: %s\n", r.UserAgent())

		mux.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return s
}
