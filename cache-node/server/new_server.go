package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/loop"
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
	Error   error           `json:"error"`
	Message string          `json:"message"`
	Value   loop.CacheEntry `json:"value"`
}

// func handleGet(eventLoop EventLoop, r *http.Request) (loop.CacheEventResponse, error) {
// 	// read body as json encoding
// 	var data RequestBody
// 	err := decodeRequestBody(r.Body, &data)
// 	if err != nil {
// 		return loop.CacheEventResponse{}, err
// 	}

// 	event, respChan, errChan := loop.CreateGetEvent(data.Key)
// 	eventLoop.Send(event)

// 	select {
// 	case resp := <-respChan:
// 		return resp, nil
// 	case err := <-errChan:
// 		return loop.CacheEventResponse{}, err
// 	}
// }

func decodeRequestBody(r io.ReadCloser, data *RequestBody) error {
	defer r.Close()
	return json.NewDecoder(r).Decode(data)
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
