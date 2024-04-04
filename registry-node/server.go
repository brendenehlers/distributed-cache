package registry

import "net/http"

type Server interface {
	HandleRegister(w http.ResponseWriter, r *http.Request)
	HandleUnregister(w http.ResponseWriter, r *http.Request)
	HandleGetNode(w http.ResponseWriter, r *http.Request)
	Start()
	Stop()
}
