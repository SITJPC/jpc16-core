package playerRepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/generic"
	"jpc16-core/type/response"
)

func FindById(id *primitive.ObjectID) (*collection.Player, *response.ErrorInstance) {
	player := new(collection.Player)
	if err := mng.PlayerCollection.FindByID(id, player); err != nil {
		return nil, response.Error(true, "Unable to find player", err)
	}
	return player, nil
}

func FindByInfo(nickname string, groupId *primitive.ObjectID) (*collection.Player, *response.ErrorInstance) {
	player := new(collection.Player)
	if err := mng.PlayerCollection.First(bson.M{
		"nickname": nickname,
		"groupId":  groupId,
	}, player); err != nil {
		return nil, response.Error(true, "Unable to find player", err)
	}
	return player, nil
}

func FindByPin(pin string) (*collection.Player, *response.ErrorInstance) {
	player := new(collection.Player)
	if err := mng.PlayerCollection.First(bson.M{
		"pin": pin,
	}, player); err != nil {
		return nil, response.Error(true, "Unable to find player", err)
	}
	return player, nil
}

func FindByTeamId(teamId *primitive.ObjectID) ([]*collection.Player, *response.ErrorInstance) {
	players := make([]*collection.Player, 0)
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{
		"teamId": teamId,
	}); err != nil {
		return nil, response.Error(true, "Unable to find players", err)
	}
	return players, nil
}

func FindManyById[T generic.ObjectID](ids []T) ([]*collection.Player, *response.ErrorInstance) {
	players := make([]*collection.Player, 0)
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}); err != nil {
		return nil, response.Error(true, "Unable to find players", err)
	}
	return players, nil
}
