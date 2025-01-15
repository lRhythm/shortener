package grpc

import "github.com/lRhythm/shortener/internal/transport"

// Server - основной объект пакета для взаимодействия.
type Server struct {
	UnimplementedShortenerServer

	cfg     transport.CfgInterface
	service transport.ServiceInterface
}

const userID = "user_id"
