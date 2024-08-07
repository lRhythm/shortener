package server

import "net/http"

type repositoryInterface interface {
	Put(key, value string)
	Get(key string) (value string)
}

type Server struct {
	mux     *http.ServeMux
	storage repositoryInterface
}
