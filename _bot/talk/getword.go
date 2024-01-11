package talk

import (
	"os"
	"strings"

	"jpc16-core/util/log"
	"jpc16-core/util/text"
)

func GetWord() string {
	// * Read file
	bytes, err := os.ReadFile("resource/wordlist.txt")
	if err != nil {
		log.Fatal("Unable to read wordlist file", err)
	}

	// * Split by newline
	lines := strings.Split(string(bytes), "\n")
	index := text.Rand.Intn(len(lines))
	return lines[index]
}
