package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
)

const (
	ERROR_MSG           = "An error has occurred"
	VALUE_FOUND_MSG     = "Value found"
	VALUE_NOT_FOUND_MSG = "Value not found"
)

type Server struct {
	http.Server
	eventLoop EventLoop
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

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	var data RequestBody
	err := decodeRequestBody(r.Body, &data)
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

func (s *Server) handleSetEvent(key string, value loop.CacheEntry) (loop.CacheEventResponse, error) {
	event, r, e := loop.CreateSetEvent(key, value)

	resp, err := s.sendEvent(event, r, e)
	if err != nil {
		return loop.CacheEventResponse{}, err
	}

	return resp, nil
}

func (s *Server) handleDeleteEvent(key string) (loop.CacheEventResponse, error) {
	event, r, e := loop.CreateDeleteEvent(key)

	resp, err := s.sendEvent(event, r, e)
	if err != nil {
		return loop.CacheEventResponse{}, err
	}

	return resp, nil
}

func (s *Server) sendEvent(
	event *loop.CacheEvent,
	respChan chan loop.CacheEventResponse,
	errChan chan error,
) (loop.CacheEventResponse, error) {
	s.eventLoop.Send(event)

	select {
	case resp := <-respChan:
		return resp, nil
	case err := <-errChan:
		return loop.CacheEventResponse{}, err
	}
}

func decodeRequestBody(r io.ReadCloser, data *RequestBody) error {
	defer r.Close()
	return json.NewDecoder(r).Decode(data)
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

func encodeResponse(resp Response) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(resp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
