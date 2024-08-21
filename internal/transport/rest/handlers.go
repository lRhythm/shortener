package rest

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (s *Server) setupHandlers() *Server {
	router := s.app.Group(s.cfg.Path())
	router.Post("/", s.createHandler)
	router.Get("/:id", s.getHandler)
	return s
}

func (s *Server) createHandler(c *fiber.Ctx) error {
	address, err := url.JoinPath(c.BaseURL(), c.OriginalURL())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	shortURL, err := s.service.CreateShortURL(string(c.BodyRaw()), address)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	return c.Status(fiber.StatusCreated).Send([]byte(shortURL))
}

func (s *Server) getHandler(c *fiber.Ctx) error {
	originalURL, err := s.service.GetShortURL(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	c.Set("Location", originalURL)
	return c.Status(fiber.StatusTemporaryRedirect).Send(nil)
}
