package main

import (
	"jpc16-core/common/config"
	"jpc16-core/common/fiber"
	"jpc16-core/common/mng"
)

func main() {
	config.Init()
	mng.Init()
	fiber.Init()
}
