package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tweet struct {
	FullText    string              `json:"full_text,omitempty" bson:"full_text"     validate:"required,min=1"`
	From        string              `json:"from"                bson:"from"`
	CreatedTime primitive.Timestamp `json:"created_time"        bson:"created_time"`
	ID          primitive.ObjectID  `json:"id,omitempty"        bson:"_id,omitempty"`
}
