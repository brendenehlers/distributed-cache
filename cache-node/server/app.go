package server

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	*http.Server
	eventLoop  EventLoop
}

func New(loop EventLoop, addr string) *Server {
	handler := http.NewServeMux()

	server := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		eventLoop: loop,
	}

	handler.HandleFunc("POST /get", server.GetHandler)
	handler.HandleFunc("POST /set", server.SetHandler)
	handler.HandleFunc("POST /delete", server.DeleteHandler)

	return server
}

func (s *Server) Run() {
	go s.eventLoop.Run()

	log.Printf("Server listening on '%v'", s.Server.Addr)
	s.Server.ListenAndServe()
}

func (s *Server) Stop() {
	s.eventLoop.Stop()
	s.Server.Shutdown(context.Background())
}
