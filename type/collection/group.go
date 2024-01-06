package collection

import (
	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type Group struct {
	mh.ModelBase `bson:"-,inline"`
	Name         *string `bson:"name,omitempty"`
	Number       *int    `bson:"number,omitempty"`
}

func (r *Group) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "groups")
	return coll
}
