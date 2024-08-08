package httptransport

import "github.com/gofiber/fiber/v2"

type serviceInterface interface {
	CreateShortURL(originalURL, host string) (shortURL string, err error)
	GetShortURL(key string) (originalURL string, err error)
}

type Server struct {
	app     *fiber.App
	cfg     cfg
	service serviceInterface
}

type cfg struct {
	Host string
}
