package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type PinBody struct {
	Pin *string `json:"pin" validate:"len=6"`
}

type PinTokenResponse struct {
	Token    *string             `json:"token"`
	PlayerId *primitive.ObjectID `json:"playerId"`
}

type PairBody struct {
	Pin *string `json:"pin" validate:"len=6"`
}

type TeamNameBody struct {
	TeamName *string `json:"teamName" validate:"required"`
}
