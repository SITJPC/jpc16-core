package operateEndpoint

import (
	"github.com/gofiber/fiber/v2"

	"jpc16-core/type/payload"
	"jpc16-core/type/response"
)

// HandleAddPlayerScore
// @ID addPlayerScore
// @Summary Add Player Score
// @Tags operate
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /operate/score/player [get]
func HandleAddPlayerScore(c *fiber.Ctx) error {
	// * Parse body
	body := new(payload.ScorePlayerAdd)
	if err := c.BodyParser(body); err != nil {
		return response.Error(true, "Unable to parse body", err)
	}

	// *
}
