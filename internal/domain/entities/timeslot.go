package entities

import "time"

// TimeSlot represents a time slot for doctor availability
type TimeSlot struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	DoctorID    string    `bson:"doctor_id" json:"doctor_id"`
	StartTime   time.Time `bson:"start_time" json:"start_time"`
	EndTime     time.Time `bson:"end_time" json:"end_time"`
	IsAvailable bool      `bson:"is_available" json:"is_available"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
