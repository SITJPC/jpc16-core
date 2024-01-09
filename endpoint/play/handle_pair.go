package playEndpoint

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"jpc16-core/common/mng"
	playerRepo "jpc16-core/repository/player"
	"jpc16-core/service/play"
	"jpc16-core/type/collection"
	"jpc16-core/type/misc"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/value"
)

// HandlePair
// @ID pairNew
// @Summary Pair New
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /play/pair [post]
func HandlePair(c *fiber.Ctx) error {
	// * Parse context
	ct := c.Locals("ct").(context.Context)

	// * Parse player claims
	p := c.Locals("p").(*jwt.Token).Claims.(*misc.PlayerClaim)

	// * Parse body
	body := new(payload.PairBody)
	if err := c.BodyParser(body); err != nil {
		return response.Error(true, "Unable to parse body", err)
	}

	// * Find player
	player := new(collection.Player)
	if err := mng.PlayerCollection.FindByID(p.Id, player); err != nil {
		return response.Error(true, "Unable to find player", err)
	}
	ct = context.WithValue(ct, "player", player)

	// * Check if player is paired
	if player.TeamId != nil {
		return response.Error(true, "Player already assigned to a team")
	}

	// * Find added pair
	paired, err := playerRepo.FindPair(player.ID)
	if err != nil {
		return err
	}

	// * Check if player is added
	for _, pair := range paired {
		if pair.Hex() == body.PlayerId.Hex() {
			return response.Error(true, "Player already added to a pair")
		}
	}

	// * Create pair
	pair := &collection.TeamPair{
		AdderId: player.ID,
		AddedId: body.PlayerId,
		Active:  value.Ptr(true),
	}
	if err := mng.TeamPairCollection.Create(pair); err != nil {
		return response.Error(true, "Unable to create pair", err)
	}

	// * Emit player state
	play.EmitPlayerState(ct)

	return c.JSON(response.Info("Successfully paired"))
}
