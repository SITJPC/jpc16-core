package collection

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type Player struct {
	mh.ModelBase `bson:"-,inline"`
	Nickname     *string             `bson:"nickname,omitempty"`
	GroupId      *primitive.ObjectID `bson:"groupId,omitempty"`
	Pin          *string             `bson:"pin,omitempty"`
}

func (r *Player) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "players")
	return coll
}
