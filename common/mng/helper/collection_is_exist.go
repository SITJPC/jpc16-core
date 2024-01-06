package mh

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

var cachedCollectionList = make(map[*mongo.Database][]string)

func IsCollectionExist(database *mongo.Database, name string) bool {
	if cachedCollectionList[database] == nil {
		if collections, err := database.ListCollectionNames(mgm.Ctx(), bson.M{}); err != nil {
			log.Fatal("Unable to get collection list", err)
		} else {
			cachedCollectionList[database] = collections
		}
	}

	return value.Contain[string](cachedCollectionList[database], name)
}
