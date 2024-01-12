package leaderboard

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func EmitLeaderboardState() {
	// * Fetch last 30 scores from database
	var scores []*collection.Score
	if err := mng.ScoreCollection.SimpleFind(
		&scores,
		bson.M{},
		&options.FindOptions{
			Sort: bson.M{
				"createdAt": -1,
			},
			Limit: value.Ptr[int64](30),
		},
	); err != nil {
		log.Error("Unable to fetch scores", err)
		return
	}

	// * Construct team ids
	teamIds := make([]*primitive.ObjectID, 0)
	for _, score := range scores {
		var exist bool
		for _, teamId := range teamIds {
			if teamId.Hex() == score.TeamId.Hex() {
				exist = true
				break
			}
		}
		if !exist {
			teamIds = append(teamIds, score.TeamId)
		}
	}

	// * Query game per team

}
