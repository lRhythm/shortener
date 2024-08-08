package httptransport

import "github.com/gofiber/fiber/v2"

func (s *Server) setupHandlers() *Server {
	s.app.Post("/", s.createHandler)
	s.app.Get("/:id", s.getHandler)
	return s
}

func (s *Server) createHandler(c *fiber.Ctx) error {
	shortURL, err := s.service.CreateShortURL(string(c.BodyRaw()), c.BaseURL())
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
