package talk

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

func Talk(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid arguments"); err != nil {
			log.Error("Unable to send message", err)
		}
		return
	}

	if args[0] == "ch" {
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Invalid arguments"); err != nil {
				log.Error("Unable to send message", err)
			}
			return
		}

		// * Update session
		if args[1] == "a" || args[1] == "b" {
			if _, err := mng.MiniGameTalkConfigCollection.UpdateOne(
				mgm.Ctx(),
				bson.M{},
				bson.M{
					"$set": bson.M{
						"team" + strings.ToUpper(args[1]) + "ChannelId": m.ChannelID,
					},
				}); err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to update config"); err != nil {
					log.Error("Unable to send message", err)
				}
			}
		}

		if args[1] == "sa" || args[1] == "sb" {
			if _, err := mng.MiniGameTalkConfigCollection.UpdateOne(
				mgm.Ctx(),
				bson.M{},
				bson.M{
					"$set": bson.M{
						"spectator" + strings.ToUpper(args[1][1:]) + "ChannelId": m.ChannelID,
					},
				}); err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to update config"); err != nil {
					log.Error("Unable to send message", err)
				}
			}
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, "Config updated"); err != nil {
			log.Error("Unable to send message", err)
		}
	}

	if args[0] == "start" {
		// * Check current running session
		session := new(collection.MiniGameTalkSession)
		if err := mng.MiniGameTalkSessionCollection.First(
			bson.M{
				"endedAt": bson.M{"$exists": false},
			},
			session); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to find session (%s)", err.Error())); err != nil {
				log.Error("Unable to send message", err)
			}
			return
		} else if err == nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Session is already running"); err != nil {
				log.Error("Unable to send message", err)
			}
			return
		}

		// * Update session
		session = &collection.MiniGameTalkSession{
			WordsA:  []*string{value.Ptr(GetWord())},
			WordsB:  []*string{value.Ptr(GetWord())},
			EndedAt: nil,
		}
		if err := mng.MiniGameTalkSessionCollection.Create(session); err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to create session"); err != nil {
				log.Error("Unable to send message", err)
			}
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, "Session started"); err != nil {
			log.Error("Unable to send message", err)
		}
	}

	if args[0] == "end" {
		// * Check current running session
		session := new(collection.MiniGameTalkSession)
		if err := mng.MiniGameTalkSessionCollection.First(
			bson.M{
				"startedAt": bson.M{"$exists": true},
				"endedAt":   bson.M{"$exists": false},
			},
			session); err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, "No session is running"); err != nil {
				log.Error("Unable to send message", err)
			}
			return
		}

		// * Update session
		if _, err := mng.MiniGameTalkSessionCollection.UpdateByID(
			mgm.Ctx(),
			session.ID,
			bson.M{
				"$set": bson.M{
					"endedAt": time.Now(),
				},
			},
		); err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Unable to update session"); err != nil {
				log.Error("Unable to send message", err)
			}
		}
		// * Query active config
		config := new(collection.MiniGameTalkConfig)
		if err := mng.MiniGameTalkConfigCollection.First(
			bson.M{},
			config,
		); err != nil {
			return
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, "Session ended"); err != nil {
			log.Error("Unable to send message", err)
		}

		embed := &discordgo.MessageEmbed{
			Title: "Session ended",
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamAChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
		if _, err := s.ChannelMessageSendEmbed(*config.TeamBChannelId, embed); err != nil {
			log.Error("Unable to send message", err)
		}
	}
}
