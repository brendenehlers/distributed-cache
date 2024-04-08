package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) registerServer() bool {
	resp, err := s.createAndSendPostRequest("/register")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	return isStatusOk(resp.StatusCode)
}

func (s *Server) unregisterServer() bool {
	resp, err := s.createAndSendPostRequest("/unregister")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	return isStatusOk(resp.StatusCode)
}

func (s *Server) createAndSendPostRequest(context string) (*http.Response, error) {
	type body struct {
		Url string `json:"url"`
	}

	jsonData, err := json.Marshal(body{
		Url: s.Addr,
	})
	if err != nil {
		return nil, err
	}

	r := bytes.NewBuffer(jsonData)

	return http.Post(s.formatRegistryUrl(context), "application/json", r)
}

func (s *Server) formatRegistryUrl(context string) string {
	return fmt.Sprintf("%s%s", s.registryUrl, context)
}

func isStatusOk(status int) bool {
	return status == http.StatusOK
}
