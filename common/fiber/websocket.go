package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	iwebsocket "jpc16-core/instance/websocket"
)

func InitWs(router fiber.Router) {
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index

	router.Get("/player", websocket.New(func(conn *websocket.Conn) {
		iwebsocket.Serve(conn)
	}))
}
