package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func ContextMiddleware(c *fiber.Ctx) error {
	ct := context.Background()
	c.Locals("ct", ct)
	return c.Next()
}
