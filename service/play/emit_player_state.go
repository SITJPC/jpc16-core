package play

import (
	"context"

	"jpc16-core/instance/websocket"
	playerRepo "jpc16-core/repository/player"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
)

func init() {
	websocket.InitialFunc = EmitPlayerState
}

func EmitPlayerState(ct context.Context) {
	// * Extract player from context
	player := ct.Value("player").(*collection.Player)

	// * Refresh player
	player, err := playerRepo.FindById(player.ID)
	if err != nil {
		log.Error("Unable to find player", err)
	}

	// * Check player assigned to a team
	if player.TeamId == nil {
		GetPlayerPair(ct)
	} else {
		GetPlayerLobby(ct)
	}
}
