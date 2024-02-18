package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Film struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Director string             `bson:"director"`
}
