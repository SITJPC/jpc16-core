package play

import (
	"context"

	"jpc16-core/instance/websocket"
	"jpc16-core/type/collection"
)

func init() {
	websocket.InitialFunc = EmitPlayerState
}

func EmitPlayerState(ct context.Context) {
	// * Extract player from context
	player := ct.Value("player").(*collection.Player)

	// * Check player assigned to a team
	if player.TeamId == nil {
		GetPlayerPair(ct)
	}
}
