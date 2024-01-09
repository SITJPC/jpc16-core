package misc

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayerClaim struct {
	Id *primitive.ObjectID `json:"id"`
}

func (T *PlayerClaim) Valid() error {
	return nil
}
