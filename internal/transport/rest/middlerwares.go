package rest

import "github.com/gofiber/fiber/v2"

const (
	cookieUserID = "userId"
)

func setUserID(c *fiber.Ctx, userID string) {
	c.Locals(cookieUserID, userID)
}

func userID(c *fiber.Ctx) string {
	return c.Locals(cookieUserID).(string)
}

func (s *Server) authenticateMiddleware(c *fiber.Ctx) error {
	ID := c.Cookies(cookieUserID)
	if err := s.service.ValidateUserID(ID); err != nil {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}
	setUserID(c, ID)
	return c.Next()
}

func (s *Server) registerMiddleware(c *fiber.Ctx) error {
	ID := c.Cookies(cookieUserID)
	if err := s.service.ValidateUserID(ID); err != nil {
		ID = s.service.GenerateUserID()
		c.Cookie(&fiber.Cookie{
			Name:  cookieUserID,
			Value: ID,
		})
	}
	setUserID(c, ID)
	return c.Next()
}
