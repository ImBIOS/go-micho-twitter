package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		Email     string             `json:"email"               bson:"email"               validate:"required,email"`
		Password  string             `json:"password,omitempty"  bson:"password"            validate:"required,min=6"`
		Token     string             `json:"token,omitempty"     bson:"-"`
		Followers []string           `json:"followers,omitempty" bson:"followers,omitempty"`
		ID        primitive.ObjectID `json:"id"                  bson:"_id,omitempty"`
	}
)
