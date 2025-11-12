package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PhoneNumber  string             `bson:"phoneNumber" json:"phoneNumber"`
	Password     string             `bson:"password" json:"-"`
	ProfileImage string             `bson:"profileImage,omitempty" json:"profileImage,omitempty"`
	IsVerified   bool               `bson:"isVerified" json:"isVerified"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Code      string             `bson:"code"`
	ExpiresAt time.Time          `bson:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt"`
}
