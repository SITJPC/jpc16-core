package operateEndpoint

import (
	"github.com/gofiber/fiber/v2"

	teamRepo "jpc16-core/repository/team"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/text"
)

// HandleAddPlayerScoreRoblox
// @ID addPlayerScoreRoblox
// @Summary Add Player Score Roblox
// @Tags operate
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /operate/score/roblox [get]
func HandleAddPlayerScoreRoblox(c *fiber.Ctx) error {
	// * Retrieve local game
	game := c.Locals("game").(*collection.Game)

	// * Parse body
	body := new(payload.ScorePlayerAddRoblox)
	if err := c.BodyParser(body); err != nil {
		return response.Error(true, "Unable to parse body", err)
	}

	// * Validate body
	if err := text.Validate(body); err != nil {
		return err
	}

	// * Find group
	team, err := teamRepo.FindByNumber(body.TeamNo)
	if err != nil {
		return err
	}

	// * Create score
	score := &collection.Score{
		TeamId: team.ID,
		GameId: game.ID,
		Score:  body.Score,
	}
	if err := mng.ScoreCollection.Create(score); err != nil {
		return response.Error(true, "Unable to create score", err)
	}

	return c.JSON(response.Info(true, "Roblox score added"))
}
