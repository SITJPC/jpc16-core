package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaderboardState struct {
	Groups []*LeaderboardStateGroup `json:"groups"`
}

type LeaderboardStateGroup struct {
	No             *int                              `json:"no"`
	GroupId        *primitive.ObjectID               `json:"groupId"`
	GroupName      *string                           `json:"groupName"`
	Score          *int64                            `json:"score"`
	ScoreBreakdown []*LeaderboardStateScoreBreakdown `json:"scoreBreakdown"`
	Participate    *int64                            `json:"participate"`
}

type LeaderboardStateScoreBreakdown struct {
	GameId     *primitive.ObjectID `json:"gameId"`
	Score      *int64              `json:"score"`
	Percentage *float64            `json:"percentage"`
}
