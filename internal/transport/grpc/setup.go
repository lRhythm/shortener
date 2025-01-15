/*
Package grpc - реализация gRPC взаимодействия с сервисом.
*/
package grpc

import (
	"errors"

	"google.golang.org/grpc"

	"github.com/lRhythm/shortener/internal/transport"
)

// New - конструктор grpc.Server.
func New(cfg transport.CfgInterface, service transport.ServiceInterface) (*grpc.Server, error) {
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}
	if service == nil {
		return nil, errors.New("service must not be nil")
	}
	s := new(Server)
	s.cfg = cfg
	s.service = service
	server := grpc.NewServer()
	RegisterShortenerServer(server, s)
	return server, nil
}
