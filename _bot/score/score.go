package score

import (
	"github.com/bwmarrin/discordgo"

	"jpc16-core/util/log"
)

func Score(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid arguments"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	gameName := args[0]

}
