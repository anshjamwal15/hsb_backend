package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Image       string               `bson:"image,omitempty" json:"image,omitempty"`
	MemberCount int                  `bson:"memberCount" json:"memberCount"`
	Members     []primitive.ObjectID `bson:"members" json:"members"`
	IsPrivate   bool                 `bson:"isPrivate" json:"isPrivate"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type GroupPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	GroupID   primitive.ObjectID `bson:"groupId" json:"groupId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content" json:"content"`
	Images    []string           `bson:"images,omitempty" json:"images,omitempty"`
	Likes     int                `bson:"likes" json:"likes"`
	Comments  int                `bson:"comments" json:"comments"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type GroupComment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PostID    primitive.ObjectID `bson:"postId" json:"postId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
