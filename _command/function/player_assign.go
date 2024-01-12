package function

import (
	"fmt"
	"math"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	teamRepo "jpc16-core/repository/team"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func PlayerAssign() {
	// * Find groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		log.Fatal("Unable to get groups", err)
	}

	// * Find players
	var players []*collection.Player
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{}); err != nil {
		log.Fatal("Unable to get players", err)
	}

	// * Construct group map
	groupMap := make(map[primitive.ObjectID]*collection.Group)
	for _, group := range groups {
		groupMap[*group.ID] = group
	}

	// * Iterate groups
	for _, group := range groups {
		// * Filter players
		var filteredPlayers []*collection.Player
		for _, player := range players {
			if player.GroupId.Hex() == group.ID.Hex() {
				filteredPlayers = append(filteredPlayers, player)
			}
		}

		// * Split half
		halfIndex := int(math.Ceil(float64(len(filteredPlayers) / 2)))
		playerLists := [][]*collection.Player{filteredPlayers[:halfIndex], filteredPlayers[halfIndex:]}

		// * Create team
		count, err := teamRepo.Count()
		if err != nil {
			log.Fatal("Failed to count teams", err)
		}

		// * Create team
		for i, playerList := range playerLists {
			name := fmt.Sprintf("%s %s", *group.Name, string(rune('A'+i)))
			team := &collection.Team{
				Name:   &name,
				Number: value.Ptr(count + int64(i) + 1),
			}
			if err := mng.TeamCollection.Create(team); err != nil {
				log.Fatal("Failed to create team", err)
			}

			// * Assign team
			for _, player := range playerList {
				if _, err := mng.PlayerCollection.UpdateOne(
					mgm.Ctx(),
					bson.M{
						"_id": player.ID,
					},
					bson.M{
						"$set": bson.M{
							"teamId": team.ID,
						},
					},
				); err != nil {
					log.Fatal("Failed to update player", err)
				}
			}
		}
	}

	log.Debug("Successfully assigned players")
}
