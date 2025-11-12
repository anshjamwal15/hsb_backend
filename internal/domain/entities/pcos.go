package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PCOSAssessment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Responses []interface{}      `bson:"responses" json:"responses"`
	Score     int                `bson:"score" json:"score"`
	Result    string             `bson:"result" json:"result"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
