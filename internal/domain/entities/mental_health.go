package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MentalHealthTest struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TestName      string             `bson:"testName" json:"testName"`
	Name          string             `bson:"name" json:"name"`
	DisplayName   string             `bson:"displayName,omitempty" json:"displayName,omitempty"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Questions     []TestQuestion     `bson:"questions,omitempty" json:"questions,omitempty"`
	Sections      []TestSection      `bson:"sections,omitempty" json:"sections,omitempty"`
	Thresholds    []TestThreshold    `bson:"thresholds" json:"thresholds"`
	CriticalItems []interface{}      `bson:"criticalItems,omitempty" json:"criticalItems,omitempty"`
	CriticalNote  string             `bson:"criticalNote,omitempty" json:"criticalNote,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type TestQuestion struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	QuestionText  string             `bson:"questionText" json:"questionText"`
	ResponseType  string             `bson:"responseType" json:"responseType"` // mcq, slider, binary
	Options       []QuestionOption   `bson:"options,omitempty" json:"options,omitempty"`
	ResponseRange *ResponseRange     `bson:"responseRange,omitempty" json:"responseRange,omitempty"`
}

type QuestionOption struct {
	Label string `bson:"label" json:"label"`
	Value int    `bson:"value" json:"value"`
}

type ResponseRange struct {
	Min int `bson:"min" json:"min"`
	Max int `bson:"max" json:"max"`
}

type TestSection struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SectionTitle string             `bson:"sectionTitle" json:"sectionTitle"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	ResponseType string             `bson:"responseType" json:"responseType"`
	Questions    []TestQuestion     `bson:"questions" json:"questions"`
	Options      []QuestionOption   `bson:"options,omitempty" json:"options,omitempty"`
}

type TestThreshold struct {
	Min            int    `bson:"min" json:"min"`
	Max            int    `bson:"max" json:"max"`
	Severity       string `bson:"severity" json:"severity"`
	Recommendation string `bson:"recommendation" json:"recommendation"`
	AlertLevel     string `bson:"alertLevel" json:"alertLevel"`
}

type Question struct {
	ID       string   `bson:"id" json:"id"`
	Text     string   `bson:"text" json:"text"`
	Type     string   `bson:"type" json:"type"`
	Options  []string `bson:"options,omitempty" json:"options,omitempty"`
}

type TestResult struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	TestName       string             `bson:"testName" json:"testName"`
	TotalScore     int                `bson:"totalScore" json:"totalScore"`
	ObtainedScore  int                `bson:"obtainedScore" json:"obtainedScore"`
	Score          int                `bson:"score,omitempty" json:"score,omitempty"`
	Level          string             `bson:"level,omitempty" json:"level,omitempty"`
	Answers       map[string]interface{} `bson:"answers,omitempty" json:"answers,omitempty"`
	Recommendation string             `bson:"recommendation,omitempty" json:"recommendation,omitempty"`
	TestDate       time.Time          `bson:"testDate" json:"testDate"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
