package collection

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	mh "jpc16-core/common/mng/helper"
)

type TeamPair struct {
	mh.ModelBase `bson:"-,inline"`
	AdderId      *primitive.ObjectID `bson:"adderId,omitempty"`
	AddedId      *primitive.ObjectID `bson:"addedId,omitempty"`
	Active       *bool               `bson:"active,omitempty"`
}
