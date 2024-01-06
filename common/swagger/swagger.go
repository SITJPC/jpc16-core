package swagger

import (
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// @title 		JPC16 Core API
// @version        1.0

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @host           localhost:3000
// @BasePath       /api

func Init(app fiber.Router) {
	app.Get("*", swagger.HandlerDefault)
}
