package leaderboardEndpoint

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/value"
)

// HandleGetState
// @ID getState
// @Summary Get State
// @Tags leaderboard
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /leaderboard/state [get]
func HandleGetState(c *fiber.Ctx) error {
	// * Fetch team lists
	var teams []*collection.Team
	if err := mng.TeamCollection.SimpleFind(&teams, bson.M{}); err != nil {
		return response.Error(true, "Unable to fetch teams", err)
	}

	var scores []*collection.Score
	if err := mng.ScoreCollection.SimpleFind(&scores, bson.M{}); err != nil {
		return response.Error(true, "Unable to fetch scores", err)
	}

	var games []*collection.Game
	if err := mng.GameCollection.SimpleFind(&games, bson.M{}); err != nil {
		return response.Error(true, "Unable to fetch games", err)
	}

	// * Construct team map
	var leaderboardStateTeams []*payload.LeaderboardStateTeam
	for _, team := range teams {
		leaderboardStateTeam := &payload.LeaderboardStateTeam{
			TeamId:     team.ID,
			TeamNumber: team.Number,
			TeamName:   team.Name,
			Games:      nil,
		}

		for _, game := range games {
			found := false
			for _, score := range scores {
				if score.GameId.Hex() == game.ID.Hex() && score.TeamId.Hex() == team.ID.Hex() {
					found = true
					leaderboardStateTeam.Games = append(leaderboardStateTeam.Games, &payload.LeaderboardStateGame{
						GameId:   game.ID,
						GameName: game.Name,
						Score:    score.Score,
					})
				}
			}
			if !found {
				leaderboardStateTeam.Games = append(leaderboardStateTeam.Games, &payload.LeaderboardStateGame{
					GameId:   game.ID,
					GameName: game.Name,
					Score:    value.Ptr[int64](0),
				})
			}
		}

		leaderboardStateTeams = append(leaderboardStateTeams, leaderboardStateTeam)
	}

	return c.JSON(response.Info(&payload.LeaderboardState{
		Teams: leaderboardStateTeams,
	}))
}
