package leaderboardEndpoint

import (
	"github.com/gofiber/fiber/v2"

	"jpc16-core/type/response"
)

// HandleGetState
// @ID getState
// @Summary Get State
// @Tags leaderboard
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /leaderboard/state [get]
func HandleGetState(c *fiber.Ctx) error {
	return c.JSON(response.Info("JPC16 Leaderboard"))
}
