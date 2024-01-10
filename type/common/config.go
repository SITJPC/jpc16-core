package common

type Config struct {
	Environment  *uint8    `yaml:"environment" validate:"gte=1,lte=2"`
	Address      *string   `yaml:"address" validate:"required"`
	Cors         []*string `yaml:"cors" validate:"required"`
	MongoUri     *string   `yaml:"mongoUri" validate:"required"`
	MongoDbName  *string   `yaml:"mongoDbName" validate:"required"`
	Secret       *string   `yaml:"secret" validate:"required"`
	DiscordToken *string   `yaml:"discordToken" validate:"required"`
}
