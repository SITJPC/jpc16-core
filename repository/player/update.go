package playerRepo

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/generic"
	"jpc16-core/type/response"
)

func UpdateTeamId[T generic.ObjectID](playerIds []T, teamId T) *response.ErrorInstance {
	filter := bson.M{
		"_id": bson.M{
			"$in": playerIds,
		},
	}
	update := bson.M{
		"$set": bson.M{
			"teamId": teamId,
		},
	}
	if _, err := mng.PlayerCollection.UpdateMany(mgm.Ctx(), filter, update); err != nil {
		return response.Error(true, "Unable to update team", err)
	}
	return nil
}
