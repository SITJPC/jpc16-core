package play

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/instance/websocket"
	playerRepo "jpc16-core/repository/player"
	teamRepo "jpc16-core/repository/team"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/util/log"
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

	// * Find groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		log.Fatal("Unable to get groups", err)
	}

	// * Construct group map
	groupMap := make(map[primitive.ObjectID]*collection.Group)
	for _, group := range groups {
		groupMap[*group.ID] = group
	}

	// * Construct team members
	var mappedTeamMembers []*payload.Player
	for _, teamMember := range teamMembers {
		if teamMember.ID.Hex() != player.ID.Hex() {
			mappedTeamMembers = append(mappedTeamMembers, &payload.Player{
				Id:        teamMember.ID,
				Nickname:  teamMember.Nickname,
				Name:      teamMember.Name,
				GroupName: groupMap[*teamMember.GroupId].Name,
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
