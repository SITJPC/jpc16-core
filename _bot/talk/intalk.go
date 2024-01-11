package talk

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hyperjumptech/beda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"jpc16-core/common/mng"
	mh "jpc16-core/common/mng/helper"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func InTalk(s *discordgo.Session, m *discordgo.MessageCreate) {
	// * Query active config
	config := new(collection.MiniGameTalkConfig)
	if err := mng.MiniGameTalkConfigCollection.First(
		bson.M{},
		config,
	); err != nil {
		return
	}

	// * Check valid channel
	var team string
	if m.ChannelID == *config.TeamAChannelId {
		team = "A"
	} else if m.ChannelID == *config.TeamBChannelId {
		team = "B"
	} else {
		return
	}

	// * Query active session
	session := new(collection.MiniGameTalkSession)
	if err := mng.MiniGameTalkSessionCollection.First(
		bson.M{},
		session,
	); err != nil {
		return
	}

	// * Parse content
	content := strings.ToLower(m.Content)

	// * Calculate word difference
	pairA := beda.NewStringDiff(content, *session.WordA)
	pairB := beda.NewStringDiff(content, *session.WordB)
	lDistA := pairA.LevenshteinDistance()
	ldDistA := pairA.DamerauLevenshteinDistance(1, 0, 1, 1)
	liDistA := pairA.DamerauLevenshteinDistance(0, 1, 1, 1)
	lDistB := pairB.LevenshteinDistance()
	ldDistB := pairB.DamerauLevenshteinDistance(1, 0, 1, 1)
	liDistB := pairB.DamerauLevenshteinDistance(0, 1, 1, 1)
	jDistA := pairA.JaroDistance()
	jDistB := pairB.JaroDistance()
	tDistA := pairA.TrigramCompare()
	tDistB := pairB.TrigramCompare()

	// * Add message
	message := &collection.MiniGameTalkMessage{
		ModelBase: mh.ModelBase{},
		SessionId: session.ID,
		Team:      &team,
		Message:   &content,
		Elapsed:   value.Ptr(time.Now().Sub(*session.StartedAt)),
		Score:     value.Ptr(int64(jDistA * 1000)),
	}
	if err := mng.MiniGameTalkMessageCollection.Create(message); err != nil {
		log.Error("Unable to create message", err)
	}

	// * Query high score
	highScoreA := new(collection.MiniGameTalkMessage)
	if err := mng.MiniGameTalkMessageCollection.First(
		bson.M{
			"sessionId": session.ID,
			"team":      "A",
		},
		highScoreA,
		&options.FindOneOptions{
			Sort: bson.M{
				"score": -1,
			},
		}); err != nil {
		log.Error("Unable to query high score", err)
	}

	// * Query high score
	highScoreB := new(collection.MiniGameTalkMessage)
	if err := mng.MiniGameTalkMessageCollection.First(
		bson.M{
			"sessionId": session.ID,
			"team":      "B",
		},
		highScoreB,
		&options.FindOneOptions{
			Sort: bson.M{
				"score": -1,
			},
		}); err != nil {
		log.Error("Unable to query high score", err)
	}

	// * Construct color
	var colorA int
	var colorB int
	if m.ChannelID == *config.TeamAChannelId {
		colorA = 0x00ff00
	} else {
		colorB = 0x00ff00
	}

	// * Send rich embed
	embed := &discordgo.MessageEmbed{
		Title: "Message A Comparison: " + content,
		Color: colorA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Delete",
				Value:  fmt.Sprintf("```%d```", ldDistA),
				Inline: true,
			},
			{
				Name:   "Insert",
				Value:  fmt.Sprintf("```%d```", liDistA),
				Inline: true,
			},
			{
				Name:   "Levenshtein Distance",
				Value:  fmt.Sprintf("```%d```", lDistA),
				Inline: false,
			},
			{
				Name:   "Jaro Distance",
				Value:  fmt.Sprintf("```%f```", jDistA),
				Inline: false,
			},
			{
				Name:   "Trigram Compare",
				Value:  fmt.Sprintf("```%f```", tDistA),
				Inline: false,
			},
		},
	}
	embed2 := &discordgo.MessageEmbed{
		Title: "Message B Comparison: " + content,
		Color: colorB,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Delete",
				Value:  fmt.Sprintf("```%d```", ldDistB),
				Inline: true,
			},
			{
				Name:   "Insert",
				Value:  fmt.Sprintf("```%d```", liDistB),
				Inline: true,
			},
			{
				Name:   "Levenshtein Distance",
				Value:  fmt.Sprintf("```%d```", lDistB),
				Inline: false,
			},
			{
				Name:   "Jaro Distance",
				Value:  fmt.Sprintf("```%f```", jDistB),
				Inline: false,
			},
			{
				Name:   "Trigram Compare",
				Value:  fmt.Sprintf("```%f```", tDistB),
				Inline: false,
			},
		},
	}
	embed3 := &discordgo.MessageEmbed{
		Title: "High Score",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("Latest Score (Team %s)", team),
				Value:  fmt.Sprintf("%d", int(jDistA*1000)),
				Inline: true,
			},
		},
	}

	if highScoreA.Score != nil {
		embed3.Fields = append(embed3.Fields, &discordgo.MessageEmbedField{
			Name:   "Team A High Score",
			Value:  fmt.Sprintf("%d", *highScoreA.Score),
			Inline: true,
		})
	}
	if highScoreB.Score != nil {
		embed3.Fields = append(embed3.Fields, &discordgo.MessageEmbedField{
			Name:   "Team B High Score",
			Value:  fmt.Sprintf("%d", *highScoreB.Score),
			Inline: true,
		})
	}

	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
		log.Error("Unable to send message", err)
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed2); err != nil {
		log.Error("Unable to send message", err)
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed3); err != nil {
		log.Error("Unable to send message", err)
	}

	// * End game
	if lDistA == 0 && m.ChannelID == *config.TeamAChannelId {
		embed := &discordgo.MessageEmbed{
			Title: "Game ended, team A won!",
			Color: 0xFFA500,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Word A",
					Value:  fmt.Sprintf("```%s```", *session.WordA),
					Inline: false,
				},
				{
					Name:   "Word B",
					Value:  fmt.Sprintf("```%s```", *session.WordB),
					Inline: false,
				},
			},
			Description: fmt.Sprintf("Your score is %d", jDistA*1000),
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamAChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamBChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
	}

	if lDistB == 0 && m.ChannelID == *config.TeamBChannelId {
		embed := &discordgo.MessageEmbed{
			Title: "Game ended, team B won!",
			Color: 0xFFA500,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Word A",
					Value:  fmt.Sprintf("```%s```", *session.WordA),
					Inline: false,
				},
				{
					Name:   "Word B",
					Value:  fmt.Sprintf("```%s```", *session.WordB),
					Inline: false,
				},
			},
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamAChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamBChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
	}
}
