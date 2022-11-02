package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Post struct {
		ID      primitive.ObjectID `json:"id"      bson:"_id,omitempty"`
		To      string             `json:"to"      bson:"to"`
		From    string             `json:"from"    bson:"from"`
		Message string             `json:"message" bson:"message"`
	}
)
