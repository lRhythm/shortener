package rest

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lRhythm/shortener/internal/models"
)

// headerLocation - установка заголовка Location.
func headerLocation(c *fiber.Ctx, location string) {
	header(c, fiber.HeaderLocation, location)
}

// headerContentTypeApplicationJSON - установка заголовку Content-Type значения application/json.
func headerContentTypeApplicationJSON(c *fiber.Ctx) {
	header(c, fiber.HeaderContentType, fiber.MIMEApplicationJSON)
}

// header - установка заголовка.
func header(c *fiber.Ctx, key, val string) {
	c.Set(key, val)
}

// badRequestResponse - установка статуса состояния HTTP 400 и формирования соответствующего тела ответа.
func badRequestResponse(c *fiber.Ctx) error {
	return errorResponse(c, fiber.StatusBadRequest)
}

// internalServerErrorResponse - установка статуса состояния HTTP 500 и формирования соответствующего тела ответа.
func internalServerErrorResponse(c *fiber.Ctx) error {
	return errorResponse(c, fiber.StatusInternalServerError)
}

// errorResponse - установка статуса состояния HTTP и формирования соответствующего тела ответа.
func errorResponse(c *fiber.Ctx, status int) error {
	return c.Status(status).Send(nil)
}

// createResponse - response создания сокращенного URL.
type createResponse struct {
	Result string `json:"result"`
}

// newCreateResponse - конструктор createResponse.
func newCreateResponse(shortURL string) *createResponse {
	return &createResponse{
		Result: shortURL,
	}
}

// createItemResponse - response создания сокращенного URL в виде коллекции элементов.
type createItemResponse struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

// createItemsResponse - элемент коллекции response'а создания сокращенного URL.
type createItemsResponse []createItemResponse

// newCreateItemsResponse - конструктор createItemsResponse с помощью models.Rows.
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
