package rest

import "github.com/gofiber/fiber/v2"

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
