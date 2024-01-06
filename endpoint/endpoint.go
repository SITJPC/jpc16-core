package endpoint

import (
	"github.com/gofiber/fiber/v2"

	"jpc16-core/endpoint/leaderboard"
)

func Init(router fiber.Router) {
	// * Leaderboard group
	leaderboard := router.Group("/leaderboard")
	leaderboard.Get("/state", leaderboardEndpoint.HandleGetState)
}
