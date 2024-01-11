package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"jpc16-core/_bot/talk"
	cc "jpc16-core/common"
	"jpc16-core/common/config"
	"jpc16-core/common/mng"
	"jpc16-core/util/log"
)

func main() {
	config.Init()
	mng.Init()

	discord, err := discordgo.New("Bot " + *cc.Config.DiscordToken)
	if err != nil {
		log.Fatal("Unable to create discord session", err)
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentGuildMessages

	if err := discord.Open(); err != nil {
		log.Fatal("Unable to open discord session", err)
	}

	log.Debug("Bot started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// * Split message into arguments
	args := strings.Split(m.Content, " ")
	if !strings.HasPrefix(args[0], "///") {
		talk.InTalk(s, m)
		return
	}

	// * Validate arguments
	if len(args) < 2 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid arguments"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Parse command
	command := strings.TrimPrefix(args[0], "///")
	switch command {
	case "talk":
		talk.Talk(s, m, args[1:])
	default:
		if _, err := s.ChannelMessageSend(m.ChannelID, "Unknown command"); err != nil {
			log.Error("Unable to send message", err)
		}
	}
}
