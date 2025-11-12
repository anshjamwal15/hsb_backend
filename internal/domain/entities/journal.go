package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Journal struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Content   string             `bson:"content" json:"content"`
	Category  string             `bson:"category" json:"category"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
