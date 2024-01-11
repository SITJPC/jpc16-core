package payload

type ScorePlayerAdd struct {
	GroupNo  *int64  `json:"groupNo" validate:"required"`
	Nickname *string `json:"nickname" validate:"required"`
	Score    *int64  `json:"score" validate:"required"`
}

type ScorePlayerAddRoblox struct {
	GroupNo *int64 `json:"groupNo" validate:"required"`
	Score   *int64 `json:"score" validate:"required"`
}
