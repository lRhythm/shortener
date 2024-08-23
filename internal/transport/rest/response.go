package rest

import "github.com/gofiber/fiber/v2"

func badRequestResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).Send(nil)
}

type createResponse struct {
	Result string `json:"result"`
}
