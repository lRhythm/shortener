package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (s *Server) setupHandlers() *Server {
	router := s.app.Group(s.cfg.Path())
	router.Post("/api/shorten", s.apiCreateHandler)
	router.Post("/", s.createHandler).Name("create")
	router.Get("/:id", s.getHandler)
	return s
}

func (s *Server) apiCreateHandler(c *fiber.Ctx) error {
	if c.Get("Content-type") != "application/json" {
		return badRequestResponse(c)
	}
	var req createRequest
	err := json.Unmarshal(c.BodyRaw(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	u, err := c.GetRouteURL("create", nil)
	if err != nil {
		return badRequestResponse(c)
	}
	a, err := url.JoinPath(c.BaseURL(), u)
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(req.URL, a)
	if err != nil {
		return badRequestResponse(c)
	}
	c.Set("Content-type", "application/json")
	return c.Status(fiber.StatusCreated).JSON(createResponse{Result: shortURL})
}

func (s *Server) createHandler(c *fiber.Ctx) error {
	a, err := url.JoinPath(c.BaseURL(), c.OriginalURL())
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(string(c.BodyRaw()), a)
	if err != nil {
		return badRequestResponse(c)
	}
	return c.Status(fiber.StatusCreated).Send([]byte(shortURL))
}

func (s *Server) getHandler(c *fiber.Ctx) error {
	originalURL, err := s.service.GetShortURL(c.Params("id"))
	if err != nil {
		return badRequestResponse(c)
	}
	c.Set("Location", originalURL)
	return c.Status(fiber.StatusTemporaryRedirect).Send(nil)
}
