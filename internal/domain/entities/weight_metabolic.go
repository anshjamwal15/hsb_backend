package entities

type WeightMetabolic struct {
	ID        string  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string  `bson:"userId" json:"userId"`
	Date      int64   `bson:"date" json:"date"`
	Weight    float64 `bson:"weight" json:"weight"`
	Height    float64 `bson:"height,omitempty" json:"height,omitempty"`
	BMI       float64 `bson:"bmi,omitempty" json:"bmi,omitempty"`
	WaistSize float64 `bson:"waistSize,omitempty" json:"waistSize,omitempty"`
	Notes     string  `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt int64   `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64   `bson:"updatedAt" json:"updatedAt"`
}
