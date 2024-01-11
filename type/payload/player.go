package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type TeamPlayer struct {
	TeamId  *primitive.ObjectID `json:"teamId"`
	Name    *string             `json:"name"`
	Number  *int64              `json:"number"`
	Players []*Player           `json:"players"`
}

type Player struct {
	Id        *primitive.ObjectID `json:"id"`
	Nickname  *string             `json:"name"`
	Name      *string             `json:"fullname"`
	GroupName *string             `json:"groupName"`
}
