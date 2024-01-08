package playerRepo

import (
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/response"
)

func FindByInfo(nickname string, groupNo int64) (*collection.Player, *response.ErrorInstance) {
	player := new(collection.Player)
	if err := mng.PlayerCollection.First(bson.M{
		"nickname": nickname,
		"groupNo":  groupNo,
	}, player); err != nil {
		return nil, response.Error(true, "Unable to find player", err)
	}
	return player, nil
}
