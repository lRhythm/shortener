package rest

import "github.com/gofiber/fiber/v2"

// cookieUserID - ключ в (*fiber.Ctx).Locals идентификатором пользователя из cookies.
const cookieUserID = "userId"

// setUserID - запись в (*fiber.Ctx).Locals идентификатора пользователя из cookies.
func setUserID(c *fiber.Ctx, userID string) {
	c.Locals(cookieUserID, userID)
}

// userID - получение из (*fiber.Ctx).Locals идентификатора пользователя из cookies.
func userID(c *fiber.Ctx) string {
	return c.Locals(cookieUserID).(string)
}

// authenticateMiddleware - middleware аутентификации пользователя по идентификатору из cookies.
func (s *Server) authenticateMiddleware(c *fiber.Ctx) error {
	ID := c.Cookies(cookieUserID)
	if err := s.service.ValidateUserID(ID); err != nil {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}
	setUserID(c, ID)
	return c.Next()
}

// registerMiddleware - middleware регистрации пользователя.
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
