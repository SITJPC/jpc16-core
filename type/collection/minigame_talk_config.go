package collection

import (
	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	mh "jpc16-core/common/mng/helper"
	"jpc16-core/util/log"
)

type MiniGameTalkConfig struct {
	mh.ModelBase        `bson:"-,inline"`
	TeamAChannelId      *string `bson:"teamAChannelId,omitempty"`
	TeamBChannelId      *string `bson:"teamBChannelId,omitempty"`
	SpectatorAChannelId *string `bson:"spectatorAChannelId,omitempty"`
	SpectatorBChannelId *string `bson:"spectatorBChannelId,omitempty"`
}

func (r *MiniGameTalkConfig) Collection() *mgm.Collection {
	coll, exist := mh.CreateCollection(c.Database, "minigame_talk_configs")
	if !exist {
		if err := coll.Create(&MiniGameTalkConfig{
			TeamAChannelId: nil,
			TeamBChannelId: nil,
		}); err != nil {
			log.Fatal("Unable to create default data for minigame talk config", err)
		}
	}
	return coll
}
