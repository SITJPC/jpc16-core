package main

import (
	"errors"

	"jpc16-core/util/log"
)

func main() {
	log.Debug("Hello, World!", "ehllo", 2111, "aaa", errors.New("aaa"))
	log.Error("Unable to run", errors.New("aaa"))
}
