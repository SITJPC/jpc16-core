package payload

type ScorePlayerAdd struct {
	TeamNo   *int64  `json:"teamNo" validate:"required"`
	Nickname *string `json:"nickname" validate:"required"`
	Score    *int64  `json:"score" validate:"required"`
}

type ScorePlayerAddRoblox struct {
	TeamNo *int64 `json:"teamNo" validate:"required"`
	Score  *int64 `json:"score" validate:"required"`
}
