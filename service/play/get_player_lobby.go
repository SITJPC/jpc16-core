package play

import (
	"context"

	"jpc16-core/instance/websocket"
	playerRepo "jpc16-core/repository/player"
	teamRepo "jpc16-core/repository/team"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func GetPlayerLobby(ct context.Context) {
	// * Extract player from context
	player := ct.Value("player").(*collection.Player)

	// * Find group
	team, err := teamRepo.FindById(player.TeamId)
	if err != nil {
		log.Error("Unable to find team", err)
		return
	}

	// * Find team members
	teamMembers, err := playerRepo.FindByTeamId(player.TeamId)
	if err != nil {
		log.Error("Unable to find team members", err)
		return
	}

	// * Construct team members
	var mappedTeamMembers []*payload.Player
	for _, teamMember := range teamMembers {
		if teamMember.ID.Hex() != player.ID.Hex() {
			mappedTeamMembers = append(mappedTeamMembers, &payload.Player{
				Id:        teamMember.ID,
				Nickname:  teamMember.Nickname,
				Name:      value.Ptr("Sirawit P."),
				GroupName: value.Ptr("Group 111"),
			})
		}
	}

	// * Construct outbound
	outbound := &websocket.OutboundMessage{
		Event: websocket.PlayerStateEvent,
		Payload: map[string]any{
			"page": "lobby",
			"profile": map[string]any{
				"nickname":   player.Nickname,
				"teamName":   team.Name,
				"teamNumber": team.Number,
			},
			"teamMembers": mappedTeamMembers,
		},
	}

	websocket.Emit(*player.ID, outbound)
}
