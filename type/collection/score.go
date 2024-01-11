package collection

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type Score struct {
	mh.ModelBase `bson:"-,inline"`
	GroupId      *primitive.ObjectID `bson:"groupId,omitempty"`
	PlayerId     *primitive.ObjectID `bson:"playerId,omitempty"`
	GameId       *primitive.ObjectID `bson:"gameId,omitempty"`
	Score        *int64              `bson:"score,omitempty"`
}

func (r *Score) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "scores")
	return coll
}
