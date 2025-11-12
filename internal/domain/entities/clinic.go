package entities

type Clinic struct {
	ID        string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string   `bson:"name" json:"name"`
	Address   string   `bson:"address" json:"address"`
	City      string   `bson:"city" json:"city"`
	Phone     string   `bson:"phone" json:"phone"`
	Services  []string `bson:"services" json:"services"`
	Rating    float64  `bson:"rating" json:"rating"`
	Image     string   `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt int64    `bson:"createdAt" json:"createdAt"`
}

type ClinicBooking struct {
	ID            string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        string `bson:"userId" json:"userId"`
	ClinicID      string `bson:"clinicId" json:"clinicId"`
	Service       string `bson:"service" json:"service"`
	Date          int64  `bson:"date" json:"date"`
	TimeSlot      string `bson:"timeSlot" json:"timeSlot"`
	Amount        int    `bson:"amount" json:"amount"`
	PaymentID     string `bson:"paymentId,omitempty" json:"paymentId,omitempty"`
	PaymentStatus string `bson:"paymentStatus" json:"paymentStatus"` // pending, completed, failed
	Status        string `bson:"status" json:"status"`               // pending, confirmed, completed, cancelled
	CreatedAt     int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt     int64  `bson:"updatedAt" json:"updatedAt"`
}
