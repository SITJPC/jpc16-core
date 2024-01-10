package function

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/enum"
	"jpc16-core/util/log"
	"jpc16-core/util/text"
	"jpc16-core/util/value"
)

func GameImport() {
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

		// * Generate token
		token := text.Random(text.RandomSet.MixedAlphaNum, 16)

		// * Validate game type
		if enum.GameType(row[1]) != enum.GameTypeAudit && enum.GameType(row[1]) != enum.GameTypeCredit {
			log.Fatal("Invalid game type", fmt.Errorf("%s is not a valid game type", row[1]))
		}

		// * Create player
		game := &collection.Game{
			Name:  value.Ptr(row[0]),
			Type:  value.Ptr(enum.GameType(row[1])),
			Token: token,
		}
		if err := mng.GameCollection.Create(game); err != nil {
			log.Fatal("Failed to create player", err)
		}
	}

	log.Debug("Successfully imported games")
}
