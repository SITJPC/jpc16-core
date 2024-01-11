package mng

import (
	"github.com/kamva/mgm/v3"

	"jpc16-core/type/collection"
)

var PlayerCollection *mgm.Collection
var GroupCollection *mgm.Collection
var GameCollection *mgm.Collection
var ScoreCollection *mgm.Collection
var TeamCollection *mgm.Collection
var TeamPairCollection *mgm.Collection
var MiniGameTalkConfigCollection *mgm.Collection
var MiniGameTalkSessionCollection *mgm.Collection
var MiniGameTalkMessageCollection *mgm.Collection

func Collection() {
	PlayerCollection = mgm.Coll(new(collection.Player))
	GroupCollection = mgm.Coll(new(collection.Group))
	TeamCollection = mgm.Coll(new(collection.Team))
	TeamPairCollection = mgm.Coll(new(collection.TeamPair))
	GameCollection = mgm.Coll(new(collection.Game))
	ScoreCollection = mgm.Coll(new(collection.Score))
	MiniGameTalkConfigCollection = mgm.Coll(new(collection.MiniGameTalkConfig))
	MiniGameTalkSessionCollection = mgm.Coll(new(collection.MiniGameTalkSession))
	MiniGameTalkMessageCollection = mgm.Coll(new(collection.MiniGameTalkMessage))

}
