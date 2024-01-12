package function

import (
	"encoding/csv"
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/text"
	"jpc16-core/util/value"
)

func PlayerImport() {
	// * Parse flags
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	file := fs.String("file", "", "File path to import")
	if err := fs.Parse(os.Args[3:]); err != nil {
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
		if err := mng.GroupCollection.First(bson.M{"name": row[0]}, group); errors.Is(err, mongo.ErrNoDocuments) {
			count, err := mng.GroupCollection.CountDocuments(mgm.Ctx(), bson.M{})
			if err != nil {
				log.Fatal("Failed to count groups", err)
			}

			group = &collection.Group{
				Name:   value.Ptr(row[0]),
				Number: value.Ptr(count + 1),
			}
			if err := mng.GroupCollection.Create(group); err != nil {
				log.Fatal("Failed to create group", err)
			}
		}

		// * Make first letter uppercase
		nickname := strings.Title(row[1])
		name := row[2]

		// * Generate pin
		var pin *string
		for {
			pin = text.Random(text.RandomSet.Num, 6)
			var player collection.Player
			if err := mng.PlayerCollection.First(bson.M{"pin": pin}, &player); errors.Is(err, mongo.ErrNoDocuments) {
				break
			} else if err != nil {
				log.Fatal("Failed to query player for pin check", err)
			}
		}

		// * Create player
		player := &collection.Player{
			Nickname: &nickname,
			Name:     &name,
			GroupId:  group.ID,
			TeamId:   nil,
			Pin:      pin,
		}
		if err := mng.PlayerCollection.Create(player); err != nil {
			log.Fatal("Failed to create player", err)
		}

		log.Debug("Created player", "nickname", nickname, "group", *group.Name)
	}

	log.Debug("Successfully imported players")
}
