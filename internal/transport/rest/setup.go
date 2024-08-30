package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"time"
)

func New(logs *logrus.Logger, cfg cfgInterface, service serviceInterface) (*Server, error) {
	if logs == nil {
		return nil, errors.New("logs must not be nil")
	}
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}
	if service == nil {
		return nil, errors.New("service must not be nil")
	}
	s := new(Server)
	s.app = newFiberApp(logs)
	s.logs = logs
	s.cfg = cfg
	s.service = service
	return s.setupHandlers(), nil
}

func (s *Server) Listen() error {
	return s.app.Listen(s.cfg.Host())
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func newFiberApp(logs *logrus.Logger) *fiber.App {
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)
	app.Use(
		logger.New(
			logger.Config{
				Format:     "{\"time\":\"${time}\", \"uri\": \"${protocol}://${host}${path}\", \"method\": \"${method}\", \"duration\": \"${latency}\", \"status\": \"${status}\", \"size\": \"${bytesSent}\"}\n",
				Output:     logs.Out,
				TimeFormat: time.DateTime,
			},
		),
	)
	app.Use(
		compress.New(
			compress.Config{
				Level: compress.LevelBestSpeed,
			},
		),
	)
	return app
}
