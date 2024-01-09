package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type PinBody struct {
	Pin *string `json:"pin" validate:"len=6"`
}

type PinTokenResponse struct {
	Token *string `json:"token"`
}

type PairBody struct {
	PlayerId *primitive.ObjectID `json:"playerId" validate:"required"`
}
