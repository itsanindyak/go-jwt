package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `json:"first_name" validate:"required,min=3,max=20"`
	LastName     string             `json:"last_name" validate:"required,min=3,max=20"`
	Password     string             `json:"password" validate:"required,min=6"`
	Email        string             `json:"email" validate:"required,email"`
	UserType     string             `json:"user_type" validate:"required,oneof=ADMIN MODERATOR USER"`
	RefreshToken string             `json:"refresh_token,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
