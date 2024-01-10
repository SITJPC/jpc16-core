package function

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"jpc16-core/common/mng"
	"jpc16-core/util/log"
)

func Clear() {
	// * Clear game collection
	if result, err := mng.GameCollection.DeleteMany(mgm.Ctx(), bson.M{}); err != nil {
		log.Error("Failed to clear game collection", err)
	} else {
		log.Debug("Cleared game collection", "count", result.DeletedCount)
	}

	// * Clear player collection
	if result, err := mng.PlayerCollection.DeleteMany(mgm.Ctx(), bson.M{}); err != nil {
		log.Error("Failed to clear player collection", err)
	} else {
		log.Debug("Cleared player collection", "count", result.DeletedCount)
	}

	// * Clear group collection
	if result, err := mng.GroupCollection.DeleteMany(mgm.Ctx(), bson.M{}); err != nil {
		log.Error("Failed to clear group collection", err)
	} else {
		log.Debug("Cleared group collection", "count", result.DeletedCount)
	}
}
