package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type serviceInterface interface {
	DBPing() (err error)
	CreateShortURL(originalURL, address string) (shortURL string, err error)
	GetShortURL(key string) (originalURL string, err error)
}

type cfgInterface interface {
	Host() string
	Path() string
}

type Server struct {
	app     *fiber.App
	logs    *logrus.Logger
	cfg     cfgInterface
	service serviceInterface
}
