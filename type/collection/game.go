package collection

import (
	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
	"jpc16-core/type/enum"
)

type Game struct {
	mh.ModelBase `bson:"-,inline"`
	Name         *string        `bson:"name,omitempty"`
	Type         *enum.GameType `bson:"type,omitempty"`
	Token        *string        `bson:"token,omitempty"`
}

func (r *Game) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "games")
	return coll
}
