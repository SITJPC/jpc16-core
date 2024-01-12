package talk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hyperjumptech/beda"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	var isSpectator bool
	if m.ChannelID == *config.TeamAChannelId || m.ChannelID == *config.SpectatorAChannelId {
		team = "A"
	} else if m.ChannelID == *config.TeamBChannelId || m.ChannelID == *config.SpectatorBChannelId {
		team = "B"
	} else if m.ChannelID == *config.TeamCChannelId || m.ChannelID == *config.SpectatorCChannelId {
		team = "C"
	} else if m.ChannelID == *config.TeamDChannelId || m.ChannelID == *config.SpectatorDChannelId {
		team = "D"
	} else {
		return
	}

	if m.ChannelID == *config.SpectatorAChannelId || m.ChannelID == *config.SpectatorBChannelId || m.ChannelID == *config.SpectatorCChannelId || m.ChannelID == *config.SpectatorDChannelId {
		isSpectator = true
	}

	channelIds := []*string{
		config.TeamAChannelId,
		config.TeamBChannelId,
		config.TeamCChannelId,
		config.TeamDChannelId,
	}

	// * Case of change word
	if isSpectator {
		if m.Content != "change" {
			return
		}
		// * Query active session
		session := new(collection.MiniGameTalkSession)
		if err := mng.MiniGameTalkSessionCollection.First(
			bson.M{
				"endedAt": bson.M{"$exists": false},
			},
			session,
			&options.FindOneOptions{
				Sort: bson.M{
					"createdAt": -1,
				},
			},
		); err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to find active session"); err != nil {
				log.Error("Unable to send message", err)
			}
		}

		// * Generate word
		var sentence bool
		if *session.Mode == "sn" {
			sentence = true
		}
		word := GetWord(sentence)
		if _, err := mng.MiniGameTalkSessionCollection.UpdateOne(
			mgm.Ctx(),
			bson.M{
				"_id": session.ID,
			},
			bson.M{
				"$push": bson.M{
					"words" + team: word,
				},
			},
		); err != nil {
			log.Error("Unable to update session", err)
		}

		// * Send message
		embed := &discordgo.MessageEmbed{
			Title:       "New word generated",
			Description: fmt.Sprintf("Your new word is **%s**", word),
			Author: &discordgo.MessageEmbedAuthor{
				Name: "Team " + team,
			},
		}
		if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Query active session
	session := new(collection.MiniGameTalkSession)
	if err := mng.MiniGameTalkSessionCollection.First(
		bson.M{},
		session,
		&options.FindOneOptions{
			Sort: bson.M{
				"createdAt": -1,
			},
		},
	); err != nil {
		return
	}

	// * Parse content
	content := strings.ToLower(m.Content)

	// * Calculate word difference
	var pair *beda.StringDiff
	if team == "A" {
		pair = beda.NewStringDiff(content, *session.WordsA[len(session.WordsA)-1])
	} else if team == "B" {
		pair = beda.NewStringDiff(content, *session.WordsB[len(session.WordsB)-1])
	} else if team == "C" {
		pair = beda.NewStringDiff(content, *session.WordsC[len(session.WordsC)-1])
	} else if team == "D" {
		pair = beda.NewStringDiff(content, *session.WordsD[len(session.WordsD)-1])
	}

	// * Algorithm alter
	lDist := pair.LevenshteinDistance()
	ldDist := pair.DamerauLevenshteinDistance(1, 0, 1, 1)
	liDist := pair.DamerauLevenshteinDistance(0, 1, 1, 1)
	jDist := pair.JaroDistance()

	// * Add message
	message := &collection.MiniGameTalkMessage{
		ModelBase: mh.ModelBase{},
		SessionId: session.ID,
		Team:      &team,
		Message:   &content,
		Elapsed:   value.Ptr(time.Now().Sub(*session.CreatedAt)),
		Score:     nil,
	}
	message.Score = value.Ptr(int64(jDist * 1000))
	if err := mng.MiniGameTalkMessageCollection.Create(message); err != nil {
		log.Error("Unable to create message", err)
	}

	// * Query high score
	var highScores []int64
	for i := 'A'; i <= 'D'; i++ {
		highScore := new(collection.MiniGameTalkMessage)
		if err := mng.MiniGameTalkMessageCollection.First(
			bson.M{
				"sessionId": session.ID,
				"team":      string(i),
			},
			highScore,
			&options.FindOneOptions{
				Sort: bson.M{
					"score": -1,
				},
			}); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			log.Error("Unable to query high score for "+string(i), err)
		}
		var highScoreValue int64
		if highScore.Score != nil {
			highScoreValue = *highScore.Score
		}
		highScores = append(highScores, highScoreValue)
	}

	// * Send rich embed
	embed := &discordgo.MessageEmbed{
		Title: "Message Comparison: ```" + content + "```",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Delete",
				Value:  fmt.Sprintf("```%d```", ldDist),
				Inline: true,
			},
			{
				Name:   "Insert",
				Value:  fmt.Sprintf("```%d```", liDist),
				Inline: true,
			},
			{
				Name:   "Levenshtein Distance",
				Value:  fmt.Sprintf("```%d```", lDist),
				Inline: false,
			},
			{
				Name:   "Jaro Distance",
				Value:  fmt.Sprintf("```%f```", jDist),
				Inline: false,
			},
		},
	}
	pairA := beda.NewStringDiff(content, *session.WordsA[len(session.WordsA)-1])
	jDistA := pairA.JaroDistance()
	pairB := beda.NewStringDiff(content, *session.WordsB[len(session.WordsB)-1])
	jDistB := pairB.JaroDistance()
	pairC := beda.NewStringDiff(content, *session.WordsC[len(session.WordsC)-1])
	jDistC := pairC.JaroDistance()
	pairD := beda.NewStringDiff(content, *session.WordsD[len(session.WordsD)-1])
	jDistD := pairD.JaroDistance()
	embed2 := &discordgo.MessageEmbed{
		Title:       "Other team message comparison",
		Color:       0x00aa00,
		Description: "Showing Jaro Distance",
		Fields:      nil,
	}
	highScoreFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Team A",
			Value:  fmt.Sprintf("```%f```", jDistA),
			Inline: true,
		},
		{
			Name:   "Team B",
			Value:  fmt.Sprintf("```%f```", jDistB),
			Inline: true,
		},
		{
			Name:   "Team C",
			Value:  fmt.Sprintf("```%f```", jDistC),
			Inline: true,
		},
		{
			Name:   "Team D",
			Value:  fmt.Sprintf("```%f```", jDistD),
			Inline: true,
		},
	}
	for i, highScore := range highScoreFields {
		if i == int(team[0]-'A') {
			continue
		}
		embed2.Fields = append(embed2.Fields, highScore)
	}
	embed3 := &discordgo.MessageEmbed{
		Title: "High Score",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("Latest Score (Team %s)", team),
				Value:  strconv.FormatInt(*message.Score, 10),
				Inline: false,
			},
		},
	}

	for i := 'A'; i <= 'D'; i++ {
		embed3.Fields = append(embed3.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Team %s High Score", string(i)),
			Value:  strconv.FormatInt(highScores[i-'A'], 10),
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
	if jDistA == 1 || jDistB == 1 || jDistC == 1 || jDistD == 1 {
		if _, err := mng.MiniGameTalkMessageCollection.UpdateOne(
			mgm.Ctx(),
			bson.M{
				"_id": message.ID,
			},
			bson.M{
				"$set": bson.M{
					"score": value.Ptr(int64(1000)),
				},
			},
		); err != nil {
			log.Error("Unable to update message", err)
		}

		embed := &discordgo.MessageEmbed{
			Title: "Game ended, team " + team + " won!",
			Color: 0xFFA500,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Team A's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsA[len(session.WordsA)-1]),
					Inline: false,
				},
				{
					Name:   "Team A's score",
					Value:  fmt.Sprintf("```%d```", highScores[0]),
					Inline: true,
				},
				{
					Name:   "Team B's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsB[len(session.WordsB)-1]),
					Inline: false,
				},
				{
					Name:  "Team B's score",
					Value: fmt.Sprintf("```%d```", highScores[1]),
				},
				{
					Name:   "Team C's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsC[len(session.WordsC)-1]),
					Inline: false,
				},
				{
					Name:   "Team C's score",
					Value:  fmt.Sprintf("```%d```", highScores[2]),
					Inline: true,
				},
				{
					Name:   "Team D's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsD[len(session.WordsD)-1]),
					Inline: false,
				},
				{
					Name:   "Team D's score",
					Value:  fmt.Sprintf("```%d```", highScores[3]),
					Inline: true,
				},
			},
		}
		for _, channelId := range channelIds {
			if _, err := s.ChannelMessageSendEmbed(*channelId, embed); err != nil {
				log.Error("Unable to send message", err)
			}
		}

		// * End session
		if _, err := mng.MiniGameTalkSessionCollection.UpdateOne(
			mgm.Ctx(),
			bson.M{
				"_id": session.ID,
			},
			bson.M{
				"$set": bson.M{
					"endedAt": time.Now(),
				},
			},
		); err != nil {
			log.Error("Unable to update session", err)
		}
	}
}
