package function

import (
	"encoding/csv"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"os"
	"time"
)

func Export() {
	// * Find groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		log.Fatal("Unable to get groups", err)
	}

	// * Find players
	var players []*collection.Player
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{}); err != nil {
		log.Fatal("Unable to get players", err)
	}

	// * Construct group map
	groupMap := make(map[primitive.ObjectID]*collection.Group)
	for _, group := range groups {
		groupMap[*group.ID] = group
	}

	now := time.Now().Format("2006-01-02-15-04-05")
	file, err := os.Create(fmt.Sprintf("local/export-%s.csv", now))
	if err != nil {
		log.Error("Failed to create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, record := range players {
		writer.Write([]string{
			*record.Nickname,
			*groupMap[*record.GroupId].Name,
		})
	}
	writer.Flush()
	log.Debug("Successfully exported", "file name", file.Name())
}
