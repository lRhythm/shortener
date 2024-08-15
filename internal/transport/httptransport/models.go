package httptransport

import "github.com/gofiber/fiber/v2"

type serviceInterface interface {
	CreateShortURL(originalURL, address string) (shortURL string, err error)
	GetShortURL(key string) (originalURL string, err error)
}

type cfgInterface interface {
	Host() string
	Path() string
}

type Server struct {
	app     *fiber.App
	cfg     cfgInterface
	service serviceInterface
}
