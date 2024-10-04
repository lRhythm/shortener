package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/shortener/internal/models"
	"github.com/sirupsen/logrus"
)

type serviceInterface interface {
	commonInterface
	URLInterface
	userInterface
}

type commonInterface interface {
	Ping() (err error)
}

type URLInterface interface {
	CreateShortURL(originalURL, address, userID string) (shortURL string, err error)
	CreateBatch(rows models.Rows, address, userID string) (models.Rows, error)
	GetOriginalURL(key string) (originalURL string, isDeleted bool, err error)
	GetUserURLs(address, userID string) (rows models.Rows, err error)
	DeleteUserURLs(keys []string, userID string)
}

type userInterface interface {
	GenerateUserID() string
	ValidateUserID(userID string) error
}

type cfgInterface interface {
	Host() string
	Path() string
	CookieKey() string
}

type Server struct {
	app     *fiber.App
	logs    *logrus.Logger
	cfg     cfgInterface
	service serviceInterface
}
