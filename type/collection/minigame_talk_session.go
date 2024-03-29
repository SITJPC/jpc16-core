package collection

import (
	"time"

	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
)

type MiniGameTalkSession struct {
	mh.ModelBase `bson:"-,inline"`
	WordsA       []*string  `bson:"wordsA,omitempty"`
	WordsB       []*string  `bson:"wordsB,omitempty"`
	WordsC       []*string  `bson:"wordsC,omitempty"`
	WordsD       []*string  `bson:"wordsD,omitempty"`
	Mode         *string    `bson:"mode,omitempty"`
	EndedAt      *time.Time `bson:"endedAt,omitempty"`
}

func (r *MiniGameTalkSession) Collection() *mgm.Collection {
	coll, _ := mh.CreateCollection(c.Database, "minigame_talk_sessions")
	return coll
}
