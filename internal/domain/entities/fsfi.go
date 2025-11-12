package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FSFIResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID       primitive.ObjectID `bson:"user" json:"user"`
	DomainScores DomainScores       `bson:"domainScores" json:"domainScores"`
	Responses    []interface{}      `bson:"responses" json:"responses"`
	TotalScore   float64            `bson:"totalScore" json:"totalScore"`
	Diagnosis    string             `bson:"diagnosis" json:"diagnosis"`
	SubmittedAt  time.Time          `bson:"submittedAt" json:"submittedAt"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DomainScores struct {
	Desire       float64 `bson:"desire" json:"desire"`
	Arousal      float64 `bson:"arousal" json:"arousal"`
	Lubrication  float64 `bson:"lubrication" json:"lubrication"`
	Orgasm       float64 `bson:"orgasm" json:"orgasm"`
	Satisfaction float64 `bson:"satisfaction" json:"satisfaction"`
	Pain         float64 `bson:"pain" json:"pain"`
}
