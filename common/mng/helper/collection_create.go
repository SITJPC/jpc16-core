package mh

import (
	"os"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"

	"jpc16-core/util/log"
)

func CreateCollection(database *mongo.Database, name string) (*mgm.Collection, bool) {
	exist := IsCollectionExist(database, name)

	if !exist {
		if err := database.CreateCollection(
			mgm.Ctx(),
			name,
		); err != nil {
			log.Debug("Unable to create collection", "column name", name, err)
			os.Exit(1)
		}
	}

	collection := &mgm.Collection{
		Collection: database.Collection(name),
	}

	return collection, exist
}
