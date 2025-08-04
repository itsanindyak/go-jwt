package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required,min=3,max=20"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required,min=3,max=20"`
	Password     string             `bson:"password" json:"password" validate:"required,min=6"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	UserType     string             `bson:"user_type" json:"user_type" validate:"required,oneof=ADMIN MODERATOR USER"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	Verified     bool               `bson:"verified" json:"verified"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
