package collection

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type MiniGameTalkMessage struct {
	mh.ModelBase `bson:"-,inline"`
	SessionId    *primitive.ObjectID `bson:"sessionId,omitempty"`
	Team         *string             `bson:"team,omitempty"`
	Message      *string             `bson:"message,omitempty"`
	Elapsed      *time.Duration      `bson:"elapsed,omitempty"`
	Score        *int64              `bson:"score,omitempty"`
}

func (r *MiniGameTalkMessage) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "minigame_talk_messages")
	return coll
}
