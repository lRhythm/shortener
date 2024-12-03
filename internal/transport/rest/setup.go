/*
Package rest - реализация стандарта REST api для HTTP взаимодействия с сервисом.
*/
package rest

import (
	"errors"
	_ "net/http/pprof"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/sirupsen/logrus"
)

// New - конструктор Server.
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
	s.app = newFiberApp(logs, cfg.CookieKey())
	s.logs = logs
	s.cfg = cfg
	s.service = service
	return s.setupHandlers(), nil
}

// Listen - прослушивание HTTP запросов сервера с указанного адреса.
func (s *Server) Listen() error {
	return s.app.Listen(s.cfg.Host())
}

// Shutdown - корректно завершает работу сервера, не прерывая активные соединения.
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// newFiberApp - конструктор fiber framework для сервера.
func newFiberApp(logs *logrus.Logger, cookieKey string) *fiber.App {
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
	app.Use(
		encryptcookie.New(
			encryptcookie.Config{
				Key: cookieKey,
			},
		),
	)
	app.Use(pprof.New())
	return app
}
