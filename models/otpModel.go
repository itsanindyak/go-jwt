package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	Otp    string             `bson:"code" json:"otp" validate:"required,min=3,max=20"`
	TTL    time.Time          `bson:"expireAt" json:"expireAt" validate:"required"` 
	Used   bool               `bson:"used" json:"used"`
}
