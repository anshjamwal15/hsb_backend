package entities

import (
	"time"
)

type Doctor struct {
	ID               string             `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string             `bson:"name" json:"name"`
	Experience       string             `bson:"experience,omitempty" json:"experience,omitempty"`
	Qualifications   string             `bson:"qualifications,omitempty" json:"qualifications,omitempty"`
	Image            string             `bson:"image,omitempty" json:"image,omitempty"`
	Specialization   string             `bson:"specialization" json:"specialization"`
	Bio              string             `bson:"bio,omitempty" json:"bio,omitempty"`
	About            string             `bson:"about,omitempty" json:"about,omitempty"`
	Rating           *float64           `bson:"rating,omitempty" json:"rating,omitempty"`
	TotalReviews     *int               `bson:"totalReviews,omitempty" json:"totalReviews,omitempty"`
	Location         string             `bson:"location,omitempty" json:"location,omitempty"`
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Email            string             `bson:"email" json:"email"`
	Languages        []string           `bson:"languages,omitempty" json:"languages,omitempty"`
	ConsultationFees *ConsultationFees  `bson:"consultationFees,omitempty" json:"consultationFees,omitempty"`
	IsAvailable      *bool              `bson:"isAvailable,omitempty" json:"isAvailable,omitempty"`
	Timing           *Timing            `bson:"timing,omitempty" json:"timing,omitempty"`
	PackageIncludes  []string           `bson:"packageIncludes,omitempty" json:"packageIncludes,omitempty"`
	IsApproved       *bool              `bson:"isApproved,omitempty" json:"isApproved,omitempty"`
	IsDeleted        *bool              `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ConsultationFees struct {
	VideoCall int `bson:"videoCall" json:"videoCall"`
	AudioCall int `bson:"audioCall" json:"audioCall"`
	InClinic  int `bson:"inClinic" json:"inClinic"`
}

type Timing struct {
	From string `bson:"from" json:"from"`
	To   string `bson:"to" json:"to"`
}
