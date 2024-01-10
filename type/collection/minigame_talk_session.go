package collection

import (
	"time"

	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type MiniGameTalkSession struct {
	mh.ModelBase `bson:"-,inline"`
	WordA        *string    `bson:"wordA,omitempty"`
	WordB        *string    `bson:"wordB,omitempty"`
	TeamAIds     []*string  `bson:"teamAIds,omitempty"`
	TeamBIds     []*string  `bson:"teamBIds,omitempty"`
	StartedAt    *time.Time `bson:"startedAt,omitempty"`
	EndedAt      *time.Time `bson:"endedAt,omitempty"`
}

func (r *MiniGameTalkSession) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "minigame_talk_sessions")
	return coll
}
