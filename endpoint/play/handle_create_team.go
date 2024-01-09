package playEndpoint

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"jpc16-core/common/mng"
	playerRepo "jpc16-core/repository/player"
	teamRepo "jpc16-core/repository/team"
	"jpc16-core/service/play"
	"jpc16-core/type/collection"
	"jpc16-core/type/misc"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/value"
)

// HandleCreateTeam
// @ID teamCreate
// @Summary Team Create
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /play/team/create [post]
func HandleCreateTeam(c *fiber.Ctx) error {
	// * Parse context
	ct := c.Locals("ct").(context.Context)

	// * Parse player claims
	p := c.Locals("p").(*jwt.Token).Claims.(*misc.PlayerClaim)

	// * Parse body
	body := new(payload.TeamNameBody)
	if err := c.BodyParser(body); err != nil {
		return response.Error(true, "Unable to parse body", err)
	}

	// * Find added pair
	paired, err := playerRepo.FindPair(p.Id)
	if err != nil {
		return err
	}

	// * Check if pair is sufficient
	if len(paired) < 3 {
		return response.Error(true, "Insufficient pair members")
	}

	// * Query players
	players, err := playerRepo.FindManyById(paired)
	if err != nil {
		return err
	}

	// * Check if player is paired
	for _, player := range players {
		if player.TeamId != nil {
			return response.Error(true, "A player already assigned to a team")
		}
	}

	// * Count number of teams
	count, err := teamRepo.Count()
	if err != nil {
		return err
	}

	// * Create team
	team := &collection.Team{
		Name:   body.TeamName,
		Number: value.Ptr(count + 1),
	}
	if err := mng.TeamCollection.Create(team); err != nil {
		return response.Error(true, "Unable to create team", err)
	}

	// * Update player
	if err := playerRepo.UpdateTeamId(paired, team.ID); err != nil {
		return err
	}

	// * Emit event
	for _, player := range players {
		ct2 := context.WithValue(ct, "player", player)
		play.EmitPlayerState(ct2)
	}

	// * Return response
	return c.JSON(response.Info("Successfully created team"))
}
