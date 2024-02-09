package main

import (
	"log"

	"github.com/brendenehlers/go-distributed-cache/cache-node/server"
)

func main() {
	s := server.NewServer()

	log.Printf("Starting server on '%s'", s.Addr)
	log.Fatal(s.ListenAndServe())
}
