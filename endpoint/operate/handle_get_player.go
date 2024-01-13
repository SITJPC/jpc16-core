package operateEndpoint

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
)

// HandleGetPlayer
// @ID getPlayer
// @Summary Get Player
// @Tags operate
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /operate/player [get]
func HandleGetPlayer(c *fiber.Ctx) error {
	// * Find teams
	var teams []*collection.Team
	if err := mng.TeamCollection.SimpleFind(&teams, bson.M{}); err != nil {
		return response.Error(true, "Unable to find teams", err)
	}

	// * Find groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		return response.Error(true, "Unable to find groups", err)
	}

	// * Construct group map
	groupMap := make(map[primitive.ObjectID]*collection.Group)
	for _, group := range groups {
		groupMap[*group.ID] = group
	}

	// * Get players
	var players []*collection.Player
	if err := mng.PlayerCollection.SimpleFind(
		&players,
		bson.M{
			"teamId": bson.M{
				"$exists": true,
			},
		},
	); err != nil {
		return response.Error(true, "Unable to get players", err)
	}

	// * Create group map
	teamMap := make(map[primitive.ObjectID]*payload.TeamPlayer)
	for _, team := range teams {
		teamMap[*team.ID] = &payload.TeamPlayer{
			TeamId:  team.ID,
			Name:    team.Name,
			Number:  team.Number,
			Players: nil,
		}
	}

	// * Iterate players
	for _, player := range players {
		teamMap[*player.TeamId].Players = append(teamMap[*player.TeamId].Players, &payload.Player{
			Id:        player.ID,
			Nickname:  player.Nickname,
			Name:      nil,
			GroupName: groupMap[*player.GroupId].Name,
		})
	}

	// * Iterate group map
	teamPlayers := make([]*payload.TeamPlayer, 0)
	for _, teamPlayer := range teamMap {
		teamPlayers = append(teamPlayers, teamPlayer)
	}

	// * Sort team player by team number
	slices.SortFunc(teamPlayers, func(a *payload.TeamPlayer, b *payload.TeamPlayer) int {
		if *a.Number < *b.Number {
			return -1
		} else if *a.Number > *b.Number {
			return 1
		}
		return 0
	})

	// * Return response
	return c.JSON(response.Info(teamPlayers))
}
