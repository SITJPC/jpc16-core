package main

import (
	"jpc16-core/common/config"
	"jpc16-core/common/fiber"
	"jpc16-core/common/mng"
	"jpc16-core/instance/websocket"
)

func main() {
	config.Init()
	mng.Init()
	websocket.Init()
	fiber.Init()
}
