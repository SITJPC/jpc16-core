package payload

type ScorePlayerAdd struct {
	Token    *string `json:"token"`
	Nickname *string `json:"nickname"`
	Score    *int64  `json:"score"`
}
