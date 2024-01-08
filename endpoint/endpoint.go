package endpoint

import (
	"github.com/gofiber/fiber/v2"

	middleware "jpc16-core/endpoint/_middleware"
	operateEndpoint "jpc16-core/endpoint/operate"

	"jpc16-core/endpoint/leaderboard"
)

func Init(router fiber.Router) {
	// * Leaderboard group
	leaderboard := router.Group("/leaderboard")
	leaderboard.Get("/state", leaderboardEndpoint.HandleGetState)

	// * Operate group
	operate := router.Group("/operate", middleware.GameMiddleware)
	operate.Get("/player", operateEndpoint.HandleGetPlayer)
	operate.Post("/score/player", operateEndpoint.HandleAddPlayerScore)
}
