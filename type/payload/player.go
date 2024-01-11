package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupPlayer struct {
	GroupId *primitive.ObjectID `json:"groupId"`
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
