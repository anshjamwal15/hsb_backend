package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diagnostics struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TestName     string             `bson:"testName" json:"testName"`
	TestPrice    float64            `bson:"testPrice" json:"testPrice"`
	TestImage    string             `bson:"testImage,omitempty" json:"testImage,omitempty"`
	HealthDataID primitive.ObjectID `bson:"healthDataId,omitempty" json:"healthDataId,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DiagnosticsUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Address     string             `bson:"address" json:"address"`
	City        string             `bson:"city" json:"city"`
	Pincode     string             `bson:"pincode" json:"pincode"`
	State       string             `bson:"state" json:"state"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Timing      Timing             `bson:"timing" json:"timing"`
	IsApproved  bool               `bson:"isApproved" json:"isApproved"`
	IsDeleted   bool               `bson:"isDeleted" json:"isDeleted"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DiagnosticsBooking struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID            primitive.ObjectID `bson:"userId" json:"userId"`
	DiagnosticsID     primitive.ObjectID `bson:"diagnosticsId" json:"diagnosticsId"`
	Date              time.Time          `bson:"date" json:"date"`
	Time              string             `bson:"time" json:"time"`
	Notes             string             `bson:"notes,omitempty" json:"notes,omitempty"`
	Status            string             `bson:"status" json:"status"`
	Amount            float64            `bson:"amount" json:"amount"`
	RazorpayOrderID   string             `bson:"razorpayOrderId,omitempty" json:"razorpayOrderId,omitempty"`
	RazorpayPaymentID string             `bson:"razorpayPaymentId,omitempty" json:"razorpayPaymentId,omitempty"`
	PaymentStatus     string             `bson:"paymentStatus" json:"paymentStatus"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}
