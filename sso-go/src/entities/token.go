package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Refreshtoken struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Token string             `json:"token,omitempty" bson:"token,omitempty"`
}

type TokenPair struct {
	Access  string
	Refresh string
}
