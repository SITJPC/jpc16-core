package teamRepo

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/response"
)

func FindById(id *primitive.ObjectID) (*collection.Team, *response.ErrorInstance) {
	team := new(collection.Team)
	if err := mng.TeamCollection.FindByID(id, team); err != nil {
		return nil, response.Error(true, "Unable to find team", err)
	}
	return team, nil
}

func FindByNumber(number *int64) (*collection.Team, *response.ErrorInstance) {
	team := new(collection.Team)
	if err := mng.TeamCollection.First(
		bson.M{"number": number},
		team,
	); err != nil {
		return nil, response.Error(true, "Unable to find team", err)
	}
	return team, nil
}

func Count() (int64, *response.ErrorInstance) {
	if count, err := mng.TeamCollection.CountDocuments(mgm.Ctx(), bson.M{}); err != nil {
		return 0, response.Error(true, "Unable to count teams", err)
	} else {
		return count, nil
	}
}
