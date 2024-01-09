package play

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/instance/websocket"
	"jpc16-core/repository/group"
	"jpc16-core/repository/player"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func GetPlayerPair(ct context.Context) {
	// * Extract player from context
	player := ct.Value("player").(*collection.Player)

	// * Find group
	group, err := groupRepo.FindGroupById(*player.GroupId)
	if err != nil {
		log.Error("Unable to find group", err)
		return
	}

	// * Get player pair
	pairs, err := playerRepo.FindPair(player.ID)
	if err != nil {
		log.Error("Unable to find player pair", err)
		return
	}

	// * Filter out player
	var filteredPairs []*primitive.ObjectID
	for _, pair := range pairs {
		if pair.Hex() != player.ID.Hex() {
			filteredPairs = append(filteredPairs, pair)
		}
	}

	// * Get player info
	playerPairs, err := playerRepo.FindManyById(pairs)
	if err != nil {
		log.Error("Unable to find player pair", err)
		return
	}

	// * Map player pair
	mappedPlayerPairs, _ := value.Iterate(playerPairs, func(player *collection.Player) (*payload.Player, *response.ErrorInstance) {
		return &payload.Player{
			Id:       player.ID,
			Nickname: player.Nickname,
		}, nil
	})

	// * Construct payload
	outbound := &websocket.OutboundMessage{
		Event: websocket.PlayerStateEvent,
		Payload: map[string]any{
			"page": "pair",
			"profile": map[string]any{
				"nickname":    player.Nickname,
				"groupNumber": group.Number,
				"groupName":   group.Name,
			},
			"teamPairs": mappedPlayerPairs,
		},
	}

	websocket.Emit(*player.ID, outbound)
}