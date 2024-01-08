package function

import (
	"encoding/csv"
	"errors"
	"flag"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
	"os"
)

func Import() {
	// * Parse flags
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	file := fs.String("file", "", "File path to import")
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatal("Failed to parse flags", err)
	}

	// * Validate flags
	if file == nil || *file == "" {
		log.Fatal("File path is required", nil)
	}

	// * Read file as csv
	fileInst, err := os.Open(*file)
	if err != nil {
		log.Fatal("Failed to read file", err)
	}
	defer fileInst.Close()

	// * Create reader
	reader := csv.NewReader(fileInst)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv", err)
	}

	for _, row := range rows {
		// * Skip if empty nickname
		if row[0] == "" {
			continue
		}

		// * Find existing group
		group := new(collection.Group)
		if err := mng.GroupCollection.First(bson.M{"name": row[1]}, group); errors.Is(err, mongo.ErrNoDocuments) {
			count, err := mng.GroupCollection.CountDocuments(mgm.Ctx(), bson.M{})
			if err != nil {
				log.Fatal("Failed to count groups", err)
			}

			group = &collection.Group{
				Name:   value.Ptr(row[1]),
				Number: value.Ptr(count + 1),
			}
			if err := mng.GroupCollection.Create(group); err != nil {
				log.Fatal("Failed to create group", err)
			}
		}

		// * Create player
		player := &collection.Player{
			Nickname: value.Ptr(row[0]),
			GroupId:  group.ID,
		}
		if err := mng.PlayerCollection.Create(player); err != nil {
			log.Fatal("Failed to create player", err)
		}
	}

	log.Debug("Successfully imported players")
}
