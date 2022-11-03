package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	Message string             `json:"message,omitempty" bson:"message,omitempty"`
	From    string             `json:"from"              bson:"from"`
	To      string             `json:"to"                bson:"to"`
	ID      primitive.ObjectID `json:"_id,omitempty"     bson:"_id,omitempty"`
}
