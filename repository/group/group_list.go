package groupRepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/response"
)

func FindGroupId(groupNo *int64) (*collection.Group, *response.ErrorInstance) {
	group := new(collection.Group)
	filter := bson.M{
		"number": groupNo,
	}
	if err := mng.GroupCollection.First(filter, group); err != nil {
		return nil, response.Error(true, "Unable to find group", err)
	}
	return group, nil
}
