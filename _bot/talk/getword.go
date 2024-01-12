package talk

import (
	"os"
	"strings"

	"jpc16-core/util/log"
	"jpc16-core/util/text"
)

func GetWord(sentence bool) string {
	// * Read file
	bytes, err := os.ReadFile("resource/wordlist.txt")
	if err != nil {
		log.Fatal("Unable to read wordlist file", err)
	}

	// * Split by newline
	lines := strings.Split(string(bytes), "\n")
	index := text.Rand.Intn(len(lines))
	if sentence {
		return lines[index]
	}

	// * Split by space
	words := strings.Split(lines[index], " ")
	var longWords []string
	for _, word := range words {
		if len(word) >= 6 {
			longWords = append(longWords, word)
		}
	}
	index = text.Rand.Intn(len(longWords))
	return longWords[index]
}
