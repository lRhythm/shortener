package rest

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/lRhythm/shortener/internal/models"
)

// setupHandlers - определения обработчиков маршрутов.
func (s *Server) setupHandlers() *Server {
	router := s.app.Group(s.cfg.Path())
	router.Get("/ping", s.PingHandler)
	router.Get("/api/internal/stats", s.trustedSubnetMiddleware, s.APIInternalStatsHandler)
	router.Get(fmt.Sprintf("/:%s", pathParamID), s.GetHandler)
	router.Get("/api/user/urls", s.authenticateMiddleware, s.APIUserUrlsGetHandler)
	router.Delete("/api/user/urls", s.authenticateMiddleware, s.APIUserUrlsDeleteHandler)
	router.Post("/api/shorten", s.registerMiddleware, s.APICreateHandler)
	router.Post("/api/shorten/batch", s.registerMiddleware, s.APICreateBatchHandler)
	router.Post("/", s.registerMiddleware, s.CreateHandler)
	return s
}

// PingHandler - обработчик маршрута ping.
func (s *Server) PingHandler(c *fiber.Ctx) error {
	if err := s.service.Ping(); err != nil {
		return internalServerErrorResponse(c)
	}
	return c.Status(fiber.StatusOK).Send(nil)
}

// APIInternalStatsHandler - обработчик маршрута получения статистики: кол-ва сокращенных URL и кол-ва пользователей.
func (s *Server) APIInternalStatsHandler(c *fiber.Ctx) error {
	countURL, err := s.service.CountURL()
	if err != nil {
		return internalServerErrorResponse(c)
	}
	countUser, err := s.service.CountUser()
	if err != nil {
		return internalServerErrorResponse(c)
	}
	return c.Status(fiber.StatusOK).JSON(newStatsResponse(countURL, countUser))
}

// APICreateHandler - обработчик маршрута создания сокращенного URL.
func (s *Server) APICreateHandler(c *fiber.Ctx) error {
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

// APICreateBatchHandler - обработчик маршрута пакетного создания сокращенного URL.
func (s *Server) APICreateBatchHandler(c *fiber.Ctx) error {
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

// CreateHandler - обработчик маршрута создания сокращенного URL.
func (s *Server) CreateHandler(c *fiber.Ctx) error {
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

// APIUserUrlsGetHandler - обработчик маршрута получения сокращенных URL пользователя.
func (s *Server) APIUserUrlsGetHandler(c *fiber.Ctx) error {
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

// APIUserUrlsDeleteHandler - обработчик маршрута удаления сокращенных URL пользователя.
func (s *Server) APIUserUrlsDeleteHandler(c *fiber.Ctx) error {
	headerContentTypeApplicationJSON(c)
	var req []string
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return badRequestResponse(c)
	}
	go s.service.DeleteUserURLs(req, userID(c))
	return c.Status(fiber.StatusAccepted).Send(nil)
}

// GetHandler - обработчик маршрута получения оригинального URL по сокращенному URL.
func (s *Server) GetHandler(c *fiber.Ctx) error {
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
