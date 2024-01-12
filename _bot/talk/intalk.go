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
	if m.ChannelID == *config.TeamAChannelId || m.ChannelID == *config.SpectatorAChannelId {
		team = "A"
	} else if m.ChannelID == *config.TeamBChannelId || m.ChannelID == *config.SpectatorBChannelId {
		team = "B"
	} else {
		return
	}

	// * Case of change word
	if m.ChannelID == *config.SpectatorAChannelId || m.ChannelID == *config.SpectatorBChannelId {
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
	pairA := beda.NewStringDiff(content, *session.WordsA[len(session.WordsA)-1])
	pairB := beda.NewStringDiff(content, *session.WordsB[len(session.WordsB)-1])
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
		Elapsed:   value.Ptr(time.Now().Sub(*session.CreatedAt)),
		Score:     nil,
	}
	if team == "A" {
		message.Score = value.Ptr(int64(jDistA * 1000))
	} else {
		message.Score = value.Ptr(int64(jDistB * 1000))
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
		}); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Error("Unable to query high score", err)
	}
	var highScoreAValue int64
	if highScoreA.Score != nil {
		highScoreAValue = *highScoreA.Score
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
	var highScoreBValue int64
	if highScoreB.Score != nil {
		highScoreBValue = *highScoreB.Score
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
				Value:  strconv.FormatInt(*message.Score, 10),
				Inline: true,
			},
		},
	}

	embed3.Fields = append(embed3.Fields, &discordgo.MessageEmbedField{
		Name:   "Team A High Score",
		Value:  fmt.Sprintf("%d", highScoreAValue),
		Inline: true,
	})
	embed3.Fields = append(embed3.Fields, &discordgo.MessageEmbedField{
		Name:   "Team B High Score",
		Value:  fmt.Sprintf("%d", highScoreBValue),
		Inline: true,
	})

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
	if (lDistA == 0 || lDistB == 0) && m.ChannelID == *config.TeamAChannelId {
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
			Title: "Game ended, team A won!",
			Color: 0xFFA500,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Team A's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsA[len(session.WordsA)-1]),
					Inline: false,
				},
				{
					Name:   "Team A's score",
					Value:  fmt.Sprintf("```%d```", 1000),
					Inline: true,
				},
				{
					Name:   "Team B's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsB[len(session.WordsB)-1]),
					Inline: false,
				},
				{
					Name:  "Team B's score",
					Value: fmt.Sprintf("```%d```", highScoreBValue),
				},
			},
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamAChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamBChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
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

	if (lDistA == 0 || lDistB == 0) && m.ChannelID == *config.TeamBChannelId {
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
			Title: "Game ended, team B won!",
			Color: 0xFFA500,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Team A's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsA[len(session.WordsA)-1]),
					Inline: false,
				},
				{
					Name:   "Team A's score",
					Value:  fmt.Sprintf("```%d```", highScoreAValue),
					Inline: true,
				},
				{
					Name:   "Team B's word",
					Value:  fmt.Sprintf("```%s```", *session.WordsB[len(session.WordsB)-1]),
					Inline: false,
				},
				{
					Name:  "Team B's score",
					Value: fmt.Sprintf("```%d```", 1000),
				},
			},
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamAChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamBChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
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
