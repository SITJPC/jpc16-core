package main

import (
	"jpc16-core/common/config"
	"jpc16-core/common/fiber"
)

func main() {
	config.Init()
	fiber.Init()
}
