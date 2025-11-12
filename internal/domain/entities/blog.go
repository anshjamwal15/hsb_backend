package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Author      string             `bson:"author" json:"author"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	Category    string             `bson:"category" json:"category"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	PublishedAt time.Time          `bson:"publishedAt" json:"publishedAt"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
