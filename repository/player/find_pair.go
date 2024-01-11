package playerRepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/response"
	"jpc16-core/util/value"
)

func FindPair(playerId *primitive.ObjectID) ([]*primitive.ObjectID, *response.ErrorInstance) {
	// * Find added pair
	var directPairs []*collection.TeamPair
	if err := mng.TeamPairCollection.SimpleFind(
		&directPairs,
		bson.M{
			"$or": bson.A{
				bson.M{
					"adderId": playerId,
				},
				bson.M{
					"addedId": playerId,
				},
			},
		}); err != nil {
		return nil, response.Error(false, "Unable to find direct pair", err)
	}

	// * Construct player map
	var playerMap = make(map[primitive.ObjectID]struct{})
	for _, pair := range directPairs {
		playerMap[*pair.AdderId] = struct{}{}
		playerMap[*pair.AddedId] = struct{}{}
	}

	// * Construct player ids
	playerIds := make([]*primitive.ObjectID, 0)
	for k := range playerMap {
		playerIds = append(playerIds, &k)
	}

	// * Find cyclic pair
	var cyclicPair []*collection.TeamPair
	if err := mng.TeamPairCollection.SimpleFind(
		&cyclicPair,
		bson.M{
			"$or": bson.A{
				bson.M{
					"adderId": bson.M{
						"$in": playerIds,
					},
				},
				bson.M{
					"addedId": bson.M{
						"$in": playerIds,
					},
				},
			},
		}); err != nil {
		return nil, response.Error(false, "Unable to find cyclic pair", err)
	}

	// * Update player map
	for _, pair := range cyclicPair {
		playerMap[*pair.AdderId] = struct{}{}
		playerMap[*pair.AddedId] = struct{}{}
	}

	// * Append player id
	playerMap[*playerId] = struct{}{}

	// * Construct player ids
	playerIds = make([]*primitive.ObjectID, 0)
	for k := range playerMap {
		playerIds = append(playerIds, value.Ptr(k))
	}

	return playerIds, nil
}
