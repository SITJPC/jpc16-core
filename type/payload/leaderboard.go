package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaderboardState struct {
	Teams []*LeaderboardStateTeam `json:"groups"`
}

type LeaderboardStateTeam struct {
	No             *int                              `json:"no"`
	TeamNumber     *int64                            `json:"teamNumber"`
	TeamName       *string                           `json:"teamName"`
	Score          *int64                            `json:"score"`
	ScoreBreakdown []*LeaderboardStateScoreBreakdown `json:"scoreBreakdown"`
	Participate    *int64                            `json:"participate"`
}

type LeaderboardStateScoreBreakdown struct {
	GameId     *primitive.ObjectID `json:"gameId"`
	Score      *int64              `json:"score"`
	Percentage *float64            `json:"percentage"`
}
