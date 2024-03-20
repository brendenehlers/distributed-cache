package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
)

const (
	ERROR_MSG           = "An error has occurred"
	VALUE_FOUND_MSG     = "Value found"
	VALUE_NOT_FOUND_MSG = "Value not found"
	VALUE_SET_MSG       = "Value set successfully"
	VALUE_DELETED_MSG   = "Value deleted successfully"
)

type Server struct {
	httpServer *http.Server
	eventLoop  EventLoop
}

type RequestBody struct {
	Key   string          `json:"key"`
	Value loop.CacheEntry `json:"value"`
}

type Response struct {
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Value   loop.CacheEntry `json:"value"`
}

func NewServer(loop EventLoop, addr string) *Server {
	handler := http.NewServeMux()

	server := &Server{
		eventLoop: loop,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	handler.HandleFunc("POST /get", server.getHandler)
	handler.HandleFunc("POST /set", server.setHandler)
	handler.HandleFunc("POST /delete", server.deleteHandler)

	return server
}

func (s *Server) Run() {
	go s.eventLoop.Run()

	log.Printf("Server listening on '%v'", s.httpServer.Addr)
	s.httpServer.ListenAndServe()
}

func (s *Server) Stop() {
	s.eventLoop.Stop()
	s.httpServer.Shutdown(context.Background())
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	data, err := decodeRequestBody(r.Body)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	cacheData, err := s.handleGetEvent(data.Key)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	buf, err := encodeResponse(createGetResponse(cacheData.Ok, cacheData.Value))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.Write(buf.Bytes())
}

func (s *Server) handleGetEvent(key string) (loop.CacheEventResponse, error) {
	event, r, e := loop.CreateGetEvent(key)

	resp, err := s.sendEvent(event, r, e)
	if err != nil {
		return loop.CacheEventResponse{}, err
	}

	return resp, nil
}

func createGetResponse(ok bool, value loop.CacheEntry) Response {
	if ok {
		return Response{
			Message: VALUE_FOUND_MSG,
			Value:   value,
		}
	} else {
		return Response{
			Message: VALUE_NOT_FOUND_MSG,
		}
	}
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {
	data, err := decodeRequestBody(r.Body)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	_, err = s.handleSetEvent(data.Key, data.Value)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	buf, err := encodeResponse(createSetResponse())
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.Write(buf.Bytes())
}

func (s *Server) handleSetEvent(key string, value loop.CacheEntry) (loop.CacheEventResponse, error) {
	event, r, e := loop.CreateSetEvent(key, value)

	resp, err := s.sendEvent(event, r, e)
	if err != nil {
		return loop.CacheEventResponse{}, err
	}

	return resp, nil
}

func createSetResponse() Response {
	return Response{
		Message: VALUE_SET_MSG,
	}
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	data, err := decodeRequestBody(r.Body)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	_, err = s.handleDeleteEvent(data.Key)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	buf, err := encodeResponse(createDeleteResponse())
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.Write(buf.Bytes())
}

func (s *Server) handleDeleteEvent(key string) (loop.CacheEventResponse, error) {
	event, r, e := loop.CreateDeleteEvent(key)

	resp, err := s.sendEvent(event, r, e)
	if err != nil {
		return loop.CacheEventResponse{}, err
	}

	return resp, nil
}

func createDeleteResponse() Response {
	return Response{
		Message: VALUE_DELETED_MSG,
	}
}

func (s *Server) sendEvent(
	event *loop.CacheEvent,
	respChan chan loop.CacheEventResponse,
	errChan chan error,
) (loop.CacheEventResponse, error) {
	log.Printf("Sending Event (%v) key: '%v'", event.Type, event.Key)
	s.eventLoop.Send(event)

	select {
	case resp := <-respChan:
		return resp, nil
	case err := <-errChan:
		return loop.CacheEventResponse{}, err
	}
}

func decodeRequestBody(r io.ReadCloser) (RequestBody, error) {
	defer r.Close()
	var data RequestBody
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	enc, _ := encodeResponse(createErrorResponse(err))
	w.Write(enc.Bytes())
}

func createErrorResponse(err error) Response {
	return Response{
		Error:   err.Error(),
		Message: ERROR_MSG,
	}
}

func encodeResponse(resp Response) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(resp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
