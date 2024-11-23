package rest

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) address(c *fiber.Ctx) (string, error) {
	return url.JoinPath(c.BaseURL(), s.cfg.Path())
}
