package groupRepo

import (
	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/generic"
	"jpc16-core/type/response"
)

func FindGroupById[T generic.ObjectID](id T) (*collection.Group, *response.ErrorInstance) {
	group := new(collection.Group)
	if err := mng.GroupCollection.FindByID(id, group); err != nil {
		return nil, response.Error(true, "Unable to find group", err)
	}
	return group, nil
}
