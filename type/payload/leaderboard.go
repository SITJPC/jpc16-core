package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaderboardState struct {
	Teams []*LeaderboardStateTeam `json:"groups"`
}

type LeaderboardStateTeam struct {
	TeamId     *primitive.ObjectID     `json:"teamId"`
	TeamNumber *int64                  `json:"teamNumber"`
	TeamName   *string                 `json:"teamName"`
	Games      []*LeaderboardStateGame `json:"games"`
}

type LeaderboardStateGame struct {
	GameId   *primitive.ObjectID `json:"gameId"`
	GameName *string             `json:"gameName"`
	Score    *int64              `json:"score"`
}
