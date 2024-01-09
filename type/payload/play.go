package payload

type PinBody struct {
	Pin *string `json:"pin" validate:"len=6"`
}

type PinTokenResponse struct {
	Token *string `json:"token"`
}

type PairBody struct {
	Pin *string `json:"pin" validate:"len=6"`
}

type TeamNameBody struct {
	TeamName *string `json:"teamName" validate:"required"`
}
