package presentation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/cache-node/domain"
)

type Server struct {
	loop domain.EventLoop
	addr string
}

type Response struct {
	Error   error             `json:"error"`
	Message string            `json:"message"`
	Value   domain.CacheEntry `json:"value"`
}

func NewServer(loop domain.EventLoop, addr string) *Server {
	return &Server{
		loop: loop,
		addr: addr,
	}
}

func (server *Server) StartServerAndLoop() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /get", func(w http.ResponseWriter, r *http.Request) {
		defer errorCatcher(w)

		if val, ok := server.getHandler(r); ok {
			writeOkayResponseWithValue(w, val)
		} else {
			writeValueNotFoundResponse(w)
		}
	})

	mux.HandleFunc("POST /set", func(w http.ResponseWriter, r *http.Request) {
		defer errorCatcher(w)

		if ok := server.setHandler(r); ok {
			writeOkayResponse(w)
		} else {
			panic("unexpected not okay response from setHandler")
		}
	})

	mux.HandleFunc("POST /delete", func(w http.ResponseWriter, r *http.Request) {
		defer errorCatcher(w)

		if ok := server.deleteHandler(r); ok {
			writeOkayResponse(w)
		} else {
			panic("unexpected not okay response from deleteHandler")
		}
	})

	logger := http.NewServeMux()
	logger.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	go server.loop.Run()
	log.Fatal(http.ListenAndServe(server.addr, logger))
}

func errorCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		writeError(w, r)
	}
}

func writeError(w http.ResponseWriter, data any) {
	resp := buildErrorResponse(data)

	w.WriteHeader(http.StatusInternalServerError)
	writeJson(w, resp)
}

func buildErrorResponse(data any) Response {
	switch data := data.(type) {
	case error:
		return Response{
			Error:   data,
			Message: data.Error(),
		}
	default:
		return Response{
			Error:   fmt.Errorf(data.(string)),
			Message: fmt.Sprint(data.(string)),
		}
	}
}

func (server *Server) getHandler(r *http.Request) (domain.CacheEntry, bool) {
	var data struct {
		Key string `json:"key"`
	}
	readRequestBody(r.Body, &data)

	event, responseChan, errorChan := domain.CreateGetEvent(data.Key)
	go server.loop.Send(event)

	select {
	case resp := <-responseChan:
		return parseEventResponse(resp)
	case err := <-errorChan:
		panic(err)
	}
}

func (server *Server) setHandler(r *http.Request) bool {
	var data struct {
		Key   string            `json:"key"`
		Value domain.CacheEntry `json:"value"`
	}
	readRequestBody(r.Body, &data)

	event, responseChan, errorChan := domain.CreateSetEvent(data.Key, data.Value)
	go server.loop.Send(event)

	select {
	case resp := <-responseChan:
		_, ok := parseEventResponse(resp)
		return ok
	case err := <-errorChan:
		panic(err)
	}
}

func (server *Server) deleteHandler(r *http.Request) bool {
	var data struct {
		Key string `json:"key"`
	}
	readRequestBody(r.Body, &data)

	event, responseChan, errorChan := domain.CreateDeleteEvent(data.Key)
	go server.loop.Send(event)

	select {
	case resp := <-responseChan:
		_, ok := parseEventResponse(resp)
		return ok
	case err := <-errorChan:
		panic(err)
	}
}

func readRequestBody(r io.ReadCloser, v any) {
	defer r.Close()
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		panic(err)
	}
}

func parseEventResponse(eventResponse domain.CacheEventResponse) (domain.CacheEntry, bool) {
	return eventResponse.Value, eventResponse.Ok
}

func writeValueNotFoundResponse(w http.ResponseWriter) {
	resp := Response{
		Message: "Value was not found",
	}

	w.WriteHeader(http.StatusNotFound)
	writeJson(w, resp)
}

func writeOkayResponse(w http.ResponseWriter) {
	resp := Response{
		Message: "Success",
	}

	writeJson(w, resp)
}

func writeOkayResponseWithValue(w http.ResponseWriter, val domain.CacheEntry) {
	resp := Response{
		Message: "Success",
		Value:   val,
	}

	writeJson(w, resp)
}

func writeJson(w io.Writer, data any) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}
