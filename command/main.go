package main

import (
	"fmt"
	"os"

	"jpc16-core/command/function"
	"jpc16-core/common/config"
	"jpc16-core/common/mng"
)

func main() {
	fmt.Println("JPC16 Core CLI")

	// * Initialize common modules
	config.Init()
	mng.Init()

	// * Check argument
	switch os.Args[1] {
	case "clear":
		function.Clear()
	case "player":
		if os.Args[2] == "import" {
			function.PlayerImport()
		}
		if os.Args[2] == "export" {
			function.PlayerExport()
		}
	}
}
