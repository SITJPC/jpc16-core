package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v3"

	cc "jpc16-core/common"
	"jpc16-core/type/misc"
	"jpc16-core/type/response"
)

func PlayerMiddleware() fiber.Handler {
	conf := jwtware.Config{
		SigningKey:  []byte(*cc.Config.Secret),
		TokenLookup: "cookie:playerToken",
		ContextKey:  "p",
		Claims:      new(misc.PlayerClaim),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.Error(false, "JWT validation failure", err)
		},
	}

	return jwtware.New(conf)
}
