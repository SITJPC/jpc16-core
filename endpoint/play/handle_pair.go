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
	pairs, err := playerRepo.FindPair(player.ID)
	if err != nil {
		return err
	}

	// * Check if pair quota is reached
	if len(pairs) >= 3 {
		return response.Error(true, "Pair quota reached")
	}

	// * Check target player condition
	targetPlayer, err := playerRepo.FindByPin(*body.Pin)
	if err != nil {
		return err
	}

	// * Check if target player is paired
	if targetPlayer.TeamId != nil {
		return response.Error(true, "Target player already assigned to a team")
	}

	// * Find new player pair
	targetPairs, err := playerRepo.FindPair(targetPlayer.ID)
	if err != nil {
		return err
	}

	// * Check merged pair
	if len(targetPairs)+len(pairs) > 3 {
		return response.Error(true, "Merged pair has more than 3 members")
	}

	// * Check if player is added
	for _, pair := range pairs {
		if pair.Hex() == targetPlayer.ID.Hex() {
			return response.Error(true, "Player already added to a pair")
		}
	}

	// * Create pair
	pair := &collection.TeamPair{
		AdderId: player.ID,
		AddedId: targetPlayer.ID,
		Active:  value.Ptr(true),
	}
	if err := mng.TeamPairCollection.Create(pair); err != nil {
		return response.Error(true, "Unable to create pair", err)
	}

	// * Emit player state
	play.EmitPlayerState(ct)

	return c.JSON(response.Info("Successfully paired"))
}
