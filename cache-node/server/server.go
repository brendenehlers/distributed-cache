package server

import (
	"log"
	"net/http"
)

func NewServer() *http.Server {
	h := http.NewServeMux()

	h.HandleFunc("/get", getHandler)
	h.HandleFunc("/set", setHandler)
	h.HandleFunc("/delete", deleteHandler)
	h.HandleFunc("/", invalidResponse)

	// log incoming requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s: %s %s\n", r.RemoteAddr, r.Method, r.URL)
		log.Printf("User Agent: %s\n", r.UserAgent())

		h.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return s
}

func validateMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Get handler"))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Set handler"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if ok := validateMethod(r, http.MethodPost); !ok {
		invalidResponse(w, r)
		return
	}

	w.Write([]byte("Delete handler"))
}

func invalidResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid request"))
}
