package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	cc "jpc16-core/common"
)

func Cors() fiber.Handler {
	origins := ""
	for i, s := range cc.Config.Cors {
		origins += *s
		if i < len(cc.Config.Cors)-1 {
			origins += ", "
		}
	}

	config := cors.Config{
		AllowOrigins:     origins,
		AllowCredentials: true,
	}

	return cors.New(config)
}
