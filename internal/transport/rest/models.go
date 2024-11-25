package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/lRhythm/shortener/internal/models"
)

// RepositoryInterface - интерфейс для имплементации сервисным слоем.
type serviceInterface interface {
	commonInterface
	URLInterface
	userInterface
}

// commonInterface - интерфейс вспомогательных методов.
type commonInterface interface {
	Ping() (err error)
}

// URLInterface - интерфейс CRUD методов работы с URL.
type URLInterface interface {
	CreateShortURL(originalURL, address, userID string) (shortURL string, err error)
	CreateBatch(rows models.Rows, address, userID string) (models.Rows, error)
	GetOriginalURL(key string) (originalURL string, isDeleted bool, err error)
	GetUserURLs(address, userID string) (rows models.Rows, err error)
	DeleteUserURLs(keys []string, userID string)
}

// userInterface - интерфейс методов работы с пользователем.
type userInterface interface {
	GenerateUserID() string
	ValidateUserID(userID string) error
}

// cfgInterface - интерфейс методов работы с config.
type cfgInterface interface {
	Host() string
	Path() string
	CookieKey() string
}

// Server - основной объект пакета для взаимодействия.
type Server struct {
	app     *fiber.App
	logs    *logrus.Logger
	cfg     cfgInterface
	service serviceInterface
}
