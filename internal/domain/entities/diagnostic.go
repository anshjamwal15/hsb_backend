package entities

type Diagnostic struct {
	ID          string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	Price       int      `bson:"price" json:"price"`
	Category    string   `bson:"category" json:"category"`
	Tests       []string `bson:"tests" json:"tests"`
	CreatedAt   int64    `bson:"createdAt" json:"createdAt"`
}

type DiagnosticBooking struct {
	ID            string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        string `bson:"userId" json:"userId"`
	DiagnosticID  string `bson:"diagnosticId" json:"diagnosticId"`
	Date          int64  `bson:"date" json:"date"`
	TimeSlot      string `bson:"timeSlot" json:"timeSlot"`
	Amount        int    `bson:"amount" json:"amount"`
	PaymentID     string `bson:"paymentId,omitempty" json:"paymentId,omitempty"`
	PaymentStatus string `bson:"paymentStatus" json:"paymentStatus"`
	Status        string `bson:"status" json:"status"`
	CreatedAt     int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt     int64  `bson:"updatedAt" json:"updatedAt"`
}
