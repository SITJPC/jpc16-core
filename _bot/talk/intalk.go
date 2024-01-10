package talk

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hyperjumptech/beda"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
)

func InTalk(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	// * Query active session
	session := new(collection.MiniGameTalkSession)
	if err := mng.MiniGameTalkSessionCollection.First(
		bson.M{},
		session,
	); err != nil {
		return
	}

	// * Calculate word difference
	pairA := beda.NewStringDiff(content, *session.WordA)
	pairB := beda.NewStringDiff(content, *session.WordB)
	lDistA := pairA.LevenshteinDistance()
	ldDistA := pairA.DamerauLevenshteinDistance(1, 0, 0, 1)
	lidDistA := pairA.DamerauLevenshteinDistance(0, 1, 0, 1)
	lDistB := pairB.LevenshteinDistance()
	ldDistB := pairB.DamerauLevenshteinDistance(1, 0, 0, 1)
	lidDistB := pairB.DamerauLevenshteinDistance(0, 1, 0, 1)
	jDistA := pairA.JaroDistance()
	jDistB := pairB.JaroDistance()
	tDistA := pairA.TrigramCompare()
	tDistB := pairB.TrigramCompare()

	// * Send rich embed
	embed := &discordgo.MessageEmbed{
		Title: "Message A Comparison: " + content,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Delete",
				Value:  fmt.Sprintf("```%d```", ldDistA),
				Inline: true,
			},
			{
				Name:   "Insert",
				Value:  fmt.Sprintf("```%d```", lidDistA),
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
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Delete",
				Value:  fmt.Sprintf("```%d```", ldDistB),
				Inline: true,
			},
			{
				Name:   "Insert",
				Value:  fmt.Sprintf("```%d```", lidDistB),
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
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
		log.Error("Unable to send message", err)
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed2); err != nil {
		log.Error("Unable to send message", err)
	}
}
