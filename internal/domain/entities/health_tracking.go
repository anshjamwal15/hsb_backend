package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeriodCycle struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	StartDate   time.Time          `bson:"startDate" json:"startDate"`
	EndDate     time.Time          `bson:"endDate,omitempty" json:"endDate,omitempty"`
	CycleLength int                `bson:"cycleLength,omitempty" json:"cycleLength,omitempty"`
	Flow        string             `bson:"flow,omitempty" json:"flow,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PregnancyTracker struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	DueDate        time.Time          `bson:"dueDate" json:"dueDate"`
	LastPeriodDate time.Time          `bson:"lastPeriodDate" json:"lastPeriodDate"`
	Weight         float64            `bson:"weight,omitempty" json:"weight,omitempty"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type SymptomsTracking struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID                 primitive.ObjectID `bson:"userId" json:"userId"`
	Type                   string             `bson:"type" json:"type"` // Period, Pregnancy
	Date                   time.Time          `bson:"date" json:"date"`
	WhatAreYouFeelingToday []string           `bson:"whatAreYouFeelingToday,omitempty" json:"whatAreYouFeelingToday,omitempty"`
	SexAndSexDrive         []string           `bson:"sexAndSexDrive,omitempty" json:"sexAndSexDrive,omitempty"`
	Mood                   []string           `bson:"mood,omitempty" json:"mood,omitempty"`
	Symptoms               []string           `bson:"symptoms,omitempty" json:"symptoms,omitempty"`
	PhysicalActivity       []string           `bson:"physicalActivity,omitempty" json:"physicalActivity,omitempty"`
	DigestionAndStool      []string           `bson:"digestionAndStool,omitempty" json:"digestionAndStool,omitempty"`
	PregnancyTest          []string           `bson:"pregnancyTest,omitempty" json:"pregnancyTest,omitempty"`
	OvulationTest          []string           `bson:"ovulationTest,omitempty" json:"ovulationTest,omitempty"`
	VaginalDischarge       []string           `bson:"vaginalDischarge,omitempty" json:"vaginalDischarge,omitempty"`
	Other                  []string           `bson:"other,omitempty" json:"other,omitempty"`
	CreatedAt              time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt              time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type WeightMetabolicWellness struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Weight    float64            `bson:"weight" json:"weight"`
	Height    float64            `bson:"height,omitempty" json:"height,omitempty"`
	BMI       float64            `bson:"bmi,omitempty" json:"bmi,omitempty"`
	Date      time.Time          `bson:"date" json:"date"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
