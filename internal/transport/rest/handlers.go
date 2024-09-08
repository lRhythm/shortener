package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (s *Server) setupHandlers() *Server {
	router := s.app.Group(s.cfg.Path())
	router.Get("/ping", s.pingHandler)
	router.Post("/api/shorten", s.apiCreateHandler)
	router.Post("/", s.createHandler)
	router.Get(fmt.Sprintf("/:%s", pathParamID), s.getHandler)
	return s
}

func (s *Server) pingHandler(c *fiber.Ctx) error {
	if err := s.service.DBPing(); err != nil {
		return internalServerErrorResponse(c)
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

func (s *Server) apiCreateHandler(c *fiber.Ctx) error {
	if c.Get(fiber.HeaderContentType) != fiber.MIMEApplicationJSON {
		return badRequestResponse(c)
	}
	var req createRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	a, err := url.JoinPath(c.BaseURL(), s.cfg.Path())
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(req.URL, a)
	if err != nil {
		return badRequestResponse(c)
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusCreated).JSON(createResponse{Result: shortURL})
}

func (s *Server) createHandler(c *fiber.Ctx) error {
	a, err := url.JoinPath(c.BaseURL(), s.cfg.Path())
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(string(c.Body()), a)
	if err != nil {
		return badRequestResponse(c)
	}
	return c.Status(fiber.StatusCreated).Send([]byte(shortURL))
}

func (s *Server) getHandler(c *fiber.Ctx) error {
	originalURL, err := s.service.GetShortURL(c.Params(pathParamID))
	if err != nil {
		return badRequestResponse(c)
	}
	c.Set(fiber.HeaderLocation, originalURL)
	return c.Status(fiber.StatusTemporaryRedirect).Send(nil)
}
