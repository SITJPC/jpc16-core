package score

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/enum"
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
	teamArg := args[1]
	teamNumber, err := strconv.ParseInt(teamArg, 10, 64)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid team"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}
	scoreArg := args[2]
	scoreNumber, err := strconv.ParseInt(scoreArg, 10, 64)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid score"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Find game by name
	game := new(collection.Game)
	if err := mng.GameCollection.First(
		bson.M{
			"name": gameName,
		},
		game,
	); err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to find game: "+gameName); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Check for audit game
	if *game.Type == enum.GameTypeAudit && scoreNumber != 1 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Audit game should have score only 1"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Find team by number
	team := new(collection.Team)
	if err := mng.TeamCollection.First(
		bson.M{
			"number": teamNumber,
		},
		team,
	); err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to find team: %d (%s)", teamNumber, err.Error())); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Check existing score
	scoreCount, err := mng.ScoreCollection.CountDocuments(
		mgm.Ctx(),
		bson.M{
			"gameId": game.ID,
			"teamId": team.ID,
		},
	)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to count score"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Check score count
	if scoreCount > 2 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Score count exceeded"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	// * Create score
	scoreObj := &collection.Score{
		GameId: game.ID,
		TeamId: team.ID,
		Score:  &scoreNumber,
	}
	if err := mng.ScoreCollection.Create(scoreObj); err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to create score"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, "Score created"); err != nil {
		log.Error("Unable to send message", err)
	}
}
