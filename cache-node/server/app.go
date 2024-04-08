package server

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	*http.Server
	eventLoop   EventLoop
	registryUrl string
}

func New(loop EventLoop, addr string, registryUrl string) *Server {
	handler := http.NewServeMux()

	server := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		eventLoop:   loop,
		registryUrl: registryUrl,
	}

	handler.HandleFunc("POST /get", server.GetHandler)
	handler.HandleFunc("POST /set", server.SetHandler)
	handler.HandleFunc("POST /delete", server.DeleteHandler)

	return server
}

func (s *Server) Run() {
	if ok := s.registerServer(); ok {
		log.Println("Successfully registered server")
	} else {
		log.Println("Unable to register server with registry node. Continuing with cache start")
	}

	go s.eventLoop.Run()

	defer func() {
		s.handleShutdown()
	}()

	log.Printf("Server listening on '%v'", s.Server.Addr)
	s.Server.ListenAndServe()
}

func (s *Server) Stop() {
	s.handleShutdown()
}

func (s *Server) handleShutdown() {
	s.eventLoop.Stop()
	s.unregisterServer()
	s.Server.Shutdown(context.Background())
}
