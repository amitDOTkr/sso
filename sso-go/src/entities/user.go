package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email,min=6,max=32"`
	Password        string             `json:"password,omitempty" bson:"password,omitempty" validate:"required,min=6,max=16"`
	Dob             time.Time          `json:"dob,omitempty" bson:"dob,omitempty"`
	ProfileImage    string             `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	IsEmailVerified bool               `json:"isEmailVerified" bson:"isEmailVerified"`
	IsActice        bool               `json:"isActive" bson:"isActive"`
	CreatedAt       time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt      time.Time          `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	// Error           Error              `json:"error,omitempty" bson:"error,omitempty"`
}
