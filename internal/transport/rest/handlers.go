package rest

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/lRhythm/shortener/internal/models"
)

func (s *Server) setupHandlers() *Server {
	router := s.app.Group(s.cfg.Path())
	router.Get("/ping", s.pingHandler)
	router.Get(fmt.Sprintf("/:%s", pathParamID), s.getHandler)
	router.Get("/api/user/urls", s.authenticateMiddleware, s.apiUserUrlsGetHandler)
	router.Delete("/api/user/urls", s.authenticateMiddleware, s.apiUserUrlsDeleteHandler)
	router.Use(s.registerMiddleware)
	router.Post("/api/shorten", s.apiCreateHandler)
	router.Post("/api/shorten/batch", s.apiCreateBatchHandler)
	router.Post("/", s.createHandler)
	return s
}

func (s *Server) pingHandler(c *fiber.Ctx) error {
	if err := s.service.Ping(); err != nil {
		return internalServerErrorResponse(c)
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

func (s *Server) apiCreateHandler(c *fiber.Ctx) error {
	headerContentTypeApplicationJSON(c)
	var req createRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	a, err := s.address(c)
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(req.OriginalURL, a, userID(c))

	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return c.Status(fiber.StatusConflict).JSON(newCreateResponse(shortURL))
		}
		return badRequestResponse(c)
	}
	return c.Status(fiber.StatusCreated).JSON(newCreateResponse(shortURL))
}

func (s *Server) apiCreateBatchHandler(c *fiber.Ctx) error {
	headerContentTypeApplicationJSON(c)
	var req createItemsRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	a, err := s.address(c)
	if err != nil {
		return badRequestResponse(c)
	}
	rows, err := s.service.CreateBatch(req.ToRows(), a, userID(c))
	if err != nil {
		return badRequestResponse(c)
	}
	return c.Status(fiber.StatusCreated).JSON(newCreateItemsResponse(rows))
}

func (s *Server) createHandler(c *fiber.Ctx) error {
	a, err := s.address(c)
	if err != nil {
		return badRequestResponse(c)
	}
	shortURL, err := s.service.CreateShortURL(string(c.Body()), a, userID(c))
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return c.Status(fiber.StatusConflict).Send([]byte(shortURL))
		}
		return badRequestResponse(c)
	}
	return c.Status(fiber.StatusCreated).Send([]byte(shortURL))
}

func (s *Server) apiUserUrlsGetHandler(c *fiber.Ctx) error {
	headerContentTypeApplicationJSON(c)
	a, err := s.address(c)
	if err != nil {
		return badRequestResponse(c)
	}
	rows, err := s.service.GetUserURLs(a, userID(c))
	if err != nil {
		return badRequestResponse(c)
	}
	if len(rows) == 0 {
		return c.Status(fiber.StatusNoContent).Send(nil)
	}
	return c.Status(fiber.StatusOK).JSON(rows)
}

func (s *Server) apiUserUrlsDeleteHandler(c *fiber.Ctx) error {
	headerContentTypeApplicationJSON(c)
	var req []string
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	go s.service.DeleteUserURLs(req, userID(c))
	return c.Status(fiber.StatusAccepted).Send(nil)
}

func (s *Server) getHandler(c *fiber.Ctx) error {
	originalURL, isDeleted, err := s.service.GetOriginalURL(c.Params(pathParamID))
	if err != nil {
		return badRequestResponse(c)
	}
	if isDeleted {
		return c.Status(fiber.StatusGone).Send(nil)
	}
	headerLocation(c, originalURL)
	return c.Status(fiber.StatusTemporaryRedirect).Send(nil)
}
