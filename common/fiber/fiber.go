package fiber

import (
	"github.com/gofiber/fiber/v2"

	cc "jpc16-core/common"
	"jpc16-core/common/fiber/middleware"
	"jpc16-core/common/swagger"
	"jpc16-core/endpoint"
	"jpc16-core/type/response"
	"jpc16-core/util/log"
	"jpc16-core/util/text"
)

func Init() {
	// Initialize fiber instance
	app := fiber.New(fiber.Config{
		AppName:       "JPC16 Core [" + text.Commit + "]",
		ErrorHandler:  ErrorHandler,
		Prefork:       false,
		StrictRouting: true,
	})

	// Register root endpoint
	app.All("/", func(c *fiber.Ctx) error {
		return c.JSON(response.Info("JPC16 API ROOT"))
	})

	// Register API endpoints
	apiGroup := app.Group("api/")
	apiGroup.Use(middleware.Recover())
	apiGroup.Use(middleware.Cors())
	endpoint.Init(apiGroup)

	// Register websocket endpoints
	wsGroup := app.Group("ws/")
	InitWs(wsGroup)

	// Register swagger endpoint
	swaggerGroup := app.Group("swagger/")
	swagger.Init(swaggerGroup)

	// Register not found endpoint
	app.Use(NotFoundHandler)

	// Startup
	err := app.Listen(*cc.Config.Address)
	if err != nil {
		log.Fatal("Unable to start fiber instance", err)
	}
}
