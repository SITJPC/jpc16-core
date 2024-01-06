package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"jpc16-core/common"
	"jpc16-core/type/common"
	"jpc16-core/util/log"
	"jpc16-core/util/text"
)

func Init() {
	// * Declare struct
	config := new(common.Config)

	// * Load configurations to struct
	yml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Unable to read configuration file", err)
	}
	if err := yaml.Unmarshal(yml, config); err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	// * Validate configurations
	if err := text.Validator.Struct(config); err != nil {
		log.Fatal("Invalid configuration file", err)
	}

	// * Set global config
	c.Config = config
}
