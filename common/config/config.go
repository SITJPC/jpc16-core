package config

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"jpc16-core/type/common"
	"jpc16-core/util/text"
)

const tag = "config"

func Init() {
	// * Declare struct
	config := new(common.Config)

	// * Load configurations to struct
	yml, err := os.ReadFile("config.yaml")
	if err != nil {

	}
	if err := yaml.Unmarshal(yml, config); err != nil {
		logrus.Fatal("UNABLE TO PARSE YAML CONFIGURATION FILE")
	}

	// * Validate configurations
	if err := text.Validator.Struct(config); err != nil {
		logrus.Fatal("INVALID CONFIGURATION: " + err.Error())
	}

	// Apply log level configuration
	logrus.SetLevel(logrus.Level(*config.LogLevel))
	spew.Config = spew.ConfigState{Indent: "  "}
}
