package mng

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kamva/mgm/v3"

	c "jpc16-core/common"
	"jpc16-core/util/log"
)

func Init() {
	// * Initialize MGM configuration
	if err := mgm.SetDefaultConfig(
		&mgm.Config{
			CtxTimeout: 5 * time.Second,
		},
		*c.Config.MongoDbName,
		options.Client().ApplyURI(*c.Config.MongoUri),
	); err != nil {
		log.Fatal("Unable to initialize MGM", err)
	}

	// * Load connection and database
	_, _, database, err := mgm.DefaultConfigs()
	if err != nil {
		log.Fatal("Unable to load MGM connection", err)
	}

	// * Assign to module
	c.Database = database

	Collection()
}
