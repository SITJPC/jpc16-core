package main

import (
	"fmt"
	"jpc16-core/command/function"
	"jpc16-core/common/config"
	"jpc16-core/common/mng"
	"os"
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
	case "import":
		function.Import()
	case "export":
		function.Export()
	}
}
