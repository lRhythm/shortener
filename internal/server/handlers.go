package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) handlerResolver(resp http.ResponseWriter, req *http.Request) {
	m := req.Method
	pathID := strings.TrimPrefix(req.URL.Path, "/")
	if m == http.MethodPost && len(pathID) == 0 {
		s.createURL(resp, req)
		return
	}
	if m == http.MethodGet && len(pathID) > 0 {
		s.getURL(resp, pathID)
		return
	}
	resp.WriteHeader(http.StatusBadRequest)
}

func (s *Server) createURL(resp http.ResponseWriter, req *http.Request) {
	b, e := io.ReadAll(req.Body)
	if e != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	u := string(b)
	_, e = url.ParseRequestURI(u)
	if e != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	k := s.genID()
	s.storage.Put(k, u)
	resp.WriteHeader(http.StatusCreated)
	_, _ = resp.Write([]byte(fmt.Sprintf("http://%s/%s", req.Host, k)))
}

func (s *Server) getURL(resp http.ResponseWriter, pathID string) {
	u := s.storage.Get(pathID)
	if len(u) == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Location", u)
	resp.WriteHeader(http.StatusTemporaryRedirect)
}
