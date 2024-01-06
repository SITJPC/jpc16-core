package common

type Config struct {
	Environment  *uint8    `yaml:"environment" validate:"gte=1,lte=2"`
	LogLevel     *uint32   `yaml:"logLevel" validate:"required"`
	Address      *string   `yaml:"address" validate:"required"`
	FrontendRoot *string   `yaml:"frontendRoot" validate:"omitempty"`
	FrontendUrl  *string   `yaml:"frontendUrl" validate:"required"`
	BackendUrl   *string   `yaml:"backendUrl" validate:"required"`
	Cors         []*string `yaml:"cors" validate:"required"`
	MongoUri     *string   `yaml:"mongoUri" validate:"required"`
	MongoDbName  *string   `yaml:"mongoDbName" validate:"required"`
	Secret       *string   `yaml:"secret" validate:"required"`
}
