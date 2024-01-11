package leaderboardEndpoint

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
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
	// * Fetch team lists
	var teams []*collection.Team
	if err := mng.TeamCollection.SimpleFind(&teams, bson.M{}); err != nil {
		return response.Error(true, "Unable to fetch teams", err)
	}

	return c.JSON(response.Info(teams))
}
