package operateEndpoint

import (
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
	// * Get groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		return response.Error(true, "Unable to get groups", err)
	}

	// * Get players
	var players []*collection.Player
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{}); err != nil {
		return response.Error(true, "Unable to get players", err)
	}

	// * Create group map
	groupMap := make(map[primitive.ObjectID]*payload.GroupPlayer)
	for _, group := range groups {
		groupMap[*group.ID] = &payload.GroupPlayer{
			GroupId: group.ID,
			Name:    group.Name,
			Number:  group.Number,
			Players: nil,
		}
	}

	// * Iterate players
	for _, player := range players {
		groupMap[*player.GroupId].Players = append(groupMap[*player.GroupId].Players, &payload.Player{
			Id:       player.ID,
			Nickname: player.Nickname,
		})
	}

	// * Iterate group map
	groupPlayers := make([]*payload.GroupPlayer, 0)
	for _, groupPlayer := range groupMap {
		groupPlayers = append(groupPlayers, groupPlayer)
	}

	// * Return response
	return c.JSON(response.Info(groupPlayers))
}
