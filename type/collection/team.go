package collection

import (
	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type Team struct {
	mh.ModelBase `bson:"-,inline"`
	Name         *string `bson:"name,omitempty"`
	Number       *int64  `bson:"number,omitempty"`
}

func (r *Team) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "teams")
	return coll
}
