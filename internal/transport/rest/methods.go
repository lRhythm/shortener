package rest

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (s *Server) address(c *fiber.Ctx) (string, error) {
	return url.JoinPath(c.BaseURL(), s.cfg.Path())
}
