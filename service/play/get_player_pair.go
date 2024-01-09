package play

import (
	"context"

	"jpc16-core/instance/websocket"
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
	group, err := playerRepo.FindGroupById(*player.GroupId)
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
				"nickname":  player.Nickname,
				"groupNo":   group.Number,
				"groupName": group.Name,
			},
			"teamPairs": mappedPlayerPairs,
		},
	}

	websocket.Emit(*player.ID, outbound)
}
