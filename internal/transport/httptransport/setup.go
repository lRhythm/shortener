package httptransport

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func New(cfg cfgInterface, service serviceInterface) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}
	if service == nil {
		return nil, errors.New("service must not be nil")
	}
	s := new(Server)
	s.app = newFiberApp()
	s.cfg = cfg
	s.service = service
	return s.setupHandlers(), nil
}

func (s *Server) Listen() error {
	return s.app.Listen(s.cfg.Host())
}

func newFiberApp() *fiber.App {
	return fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)
}
