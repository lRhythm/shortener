package server

import "net/http"

func New(repository repositoryInterface) *Server {
	s := new(Server)
	s.mux = http.NewServeMux()
	s.storage = repository
	return s.setupHandlers()
}

func (s *Server) setupHandlers() *Server {
	s.mux.HandleFunc(`/`, s.handlerResolver)
	return s
}

func (s *Server) Listen() error {
	return http.ListenAndServe(`localhost:8080`, s.mux)
}
