package generic

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectID interface {
	primitive.ObjectID | *primitive.ObjectID
}
