package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/brendenehlers/go-distributed-cache/registry-node"
)

type HttpServer struct {
	http.Server

	registry registry.Registry
}

type RequestBody struct {
	Url string `json:"url"`
}

type ErrResponseBody struct {
	Error string `json:"error"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

type NodeResponseBody struct {
	ResponseBody
	Url string `json:"url"`
}

func New(host string, reg registry.Registry) *HttpServer {
	handler := http.NewServeMux()

	server := &HttpServer{
		Server: http.Server{
			Handler: handler,
			Addr:    host,
		},
		registry: reg,
	}

	handler.HandleFunc("POST /register", logRequest(server.HandleRegister))
	handler.HandleFunc("POST /unregister", logRequest(server.HandleUnregister))
	handler.HandleFunc("GET /node", logRequest(server.HandleGetNode))

	return server
}

func (hs *HttpServer) HandleRegister(w http.ResponseWriter, r *http.Request) {
	url, err := getUrlFromBody(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err = hs.registry.Register(url)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}

	resp := &ResponseBody{Message: "Success"}
	encodeResponse(w, resp)
}

func (hs *HttpServer) HandleUnregister(w http.ResponseWriter, r *http.Request) {
	url, err := getUrlFromBody(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err = hs.registry.Unregister(url)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	resp := &ResponseBody{Message: "Success"}
	encodeResponse(w, resp)
}

func (hs *HttpServer) HandleGetNode(w http.ResponseWriter, r *http.Request) {
	node, err := hs.registry.GetNode()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	resp := &NodeResponseBody{
		ResponseBody: ResponseBody{
			Message: "Success",
		},
		Url: node.Url,
	}
	encodeResponse(w, resp)
}

func (hs *HttpServer) Start() {
	log.Printf("Listening on '%s'", hs.Addr)
	hs.ListenAndServe()
}

func (hs *HttpServer) Stop() {
	log.Printf("Stopping server...")
	hs.Shutdown(context.Background())
}

func getUrlFromBody(r io.Reader) (string, error) {
	var body RequestBody
	err := json.NewDecoder(r).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.Url, nil
}

func handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&ErrResponseBody{Error: err.Error()}); err != nil {
		log.Println(err)
	}
}

func encodeResponse(w http.ResponseWriter, val any) {
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Println(err.Error())
	}
}
