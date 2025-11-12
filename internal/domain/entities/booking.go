package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID            primitive.ObjectID `bson:"userId" json:"userId"`
	DoctorID          primitive.ObjectID `bson:"doctorId" json:"doctorId"`
	SessionType       string             `bson:"sessionType" json:"sessionType"` // Video Call, Audio Call, In-Clinic, Chat
	Date              time.Time          `bson:"date" json:"date"`
	TimeSlot          string             `bson:"timeSlot" json:"timeSlot"`
	Status            string             `bson:"status" json:"status"` // pending, confirmed, completed, cancelled
	Amount            int                `bson:"amount" json:"amount"`
	RazorpayOrderID   string             `bson:"razorpayOrderId,omitempty" json:"razorpayOrderId,omitempty"`
	RazorpayPaymentID string             `bson:"razorpayPaymentId,omitempty" json:"razorpayPaymentId,omitempty"`
	PaymentStatus     string             `bson:"paymentStatus" json:"paymentStatus"` // pending, paid, failed
	Notes             string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}
