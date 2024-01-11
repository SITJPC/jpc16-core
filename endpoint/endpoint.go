package endpoint

import (
	"github.com/gofiber/fiber/v2"

	middleware "jpc16-core/endpoint/_middleware"
	operateEndpoint "jpc16-core/endpoint/operate"
	playEndpoint "jpc16-core/endpoint/play"

	"jpc16-core/endpoint/leaderboard"
)

func Init(router fiber.Router) {
	// * Middleware
	router.Use(middleware.ContextMiddleware)

	// * Leaderboard group
	leaderboard := router.Group("/leaderboard")
	leaderboard.Get("/state", leaderboardEndpoint.HandleGetState)

	// * Operate group
	operate := router.Group("/operate", middleware.GameMiddleware)
	operate.Get("/player", operateEndpoint.HandleGetPlayer)
	operate.Post("/score/player", operateEndpoint.HandleAddPlayerScore)
	operate.Post("/score/roblox", operateEndpoint.HandleAddPlayerScoreRoblox)

	// * Play group
	play := router.Group("/play")
	play.Post("/pin", playEndpoint.HandleEnterPin)
	play.Post("/pair", middleware.PlayerMiddleware(), playEndpoint.HandlePair)
	play.Post("/team/create", middleware.PlayerMiddleware(), playEndpoint.HandleCreateTeam)
}
