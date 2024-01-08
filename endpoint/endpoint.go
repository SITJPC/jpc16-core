package endpoint

import (
	"github.com/gofiber/fiber/v2"
	operateEndpoint "jpc16-core/endpoint/operate"

	"jpc16-core/endpoint/leaderboard"
)

func Init(router fiber.Router) {
	// * Leaderboard group
	leaderboard := router.Group("/leaderboard")
	leaderboard.Get("/state", leaderboardEndpoint.HandleGetState)

	// * Operate group
	operate := router.Group("/operate")
	operate.Get("/player", operateEndpoint.HandleGetPlayer)
}
