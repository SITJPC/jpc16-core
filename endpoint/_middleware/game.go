package middleware

import "github.com/gofiber/fiber/v2"

func GameMiddleware(c *fiber.Ctx) error {

	return c.Next()
}
