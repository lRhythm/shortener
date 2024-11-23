package rest

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lRhythm/shortener/internal/models"
)

func headerLocation(c *fiber.Ctx, location string) {
	header(c, fiber.HeaderLocation, location)
}

func headerContentTypeApplicationJSON(c *fiber.Ctx) {
	header(c, fiber.HeaderContentType, fiber.MIMEApplicationJSON)
}

func header(c *fiber.Ctx, key, val string) {
	c.Set(key, val)
}

func badRequestResponse(c *fiber.Ctx) error {
	return errorResponse(c, fiber.StatusBadRequest)
}

func internalServerErrorResponse(c *fiber.Ctx) error {
	return errorResponse(c, fiber.StatusInternalServerError)
}

func errorResponse(c *fiber.Ctx, status int) error {
	return c.Status(status).Send(nil)
}

type createResponse struct {
	Result string `json:"result"`
}

func newCreateResponse(shortURL string) *createResponse {
	return &createResponse{
		Result: shortURL,
	}
}

type createItemsResponse []createItemResponse

type createItemResponse struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

func newCreateItemsResponse(rows models.Rows) createItemsResponse {
	items := make(createItemsResponse, 0)
	for _, row := range rows {
		items = append(items, createItemResponse{
			ShortURL:      row.ShortURL,
			CorrelationID: row.CorrelationID,
		})
	}
	return items
}
