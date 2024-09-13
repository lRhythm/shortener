package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/shortener/internal/models"
	"github.com/sirupsen/logrus"
)

type serviceInterface interface {
	Ping() (err error)
	CreateShortURL(originalURL, address string) (shortURL string, err error)
	CreateBatch(rows models.Rows, address string) (models.Rows, error)
	GetOriginalURL(key string) (originalURL string, err error)
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
