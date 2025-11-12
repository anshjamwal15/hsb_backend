package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Constants for test data
const (
	TestUser1ID = "5f8d04b3b54764421b7156d0"
	TestUser2ID = "5f8d04b3b54764421b7156d1"
	TestUser3ID = "5f8d04b3b54764421b7156d2"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PhoneNumber  string             `bson:"phoneNumber" json:"phoneNumber"`
	Password     string             `bson:"password" json:"-"`
	ProfileImage string             `bson:"profileImage,omitempty" json:"profileImage,omitempty"`
	DateOfBirth  time.Time          `bson:"dateOfBirth,omitempty" json:"dateOfBirth,omitempty"`
	Gender       string             `bson:"gender,omitempty" json:"gender,omitempty"`
	IsVerified   bool               `bson:"isVerified" json:"isVerified"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Doctor struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID   `bson:"userId" json:"userId"`
	Specialties     []string             `bson:"specialties" json:"specialties"`
	Qualifications  string               `bson:"qualifications" json:"qualifications"`
	Experience      int                  `bson:"experience" json:"experience"`
	Bio             string               `bson:"bio,omitempty" json:"bio,omitempty"`
	IsAvailable     bool                 `bson:"isAvailable" json:"isAvailable"`
	ConsultationFee int                  `bson:"consultationFee" json:"consultationFee"`
	Languages       []string             `bson:"languages,omitempty" json:"languages,omitempty"`
	Education       []Education          `bson:"education,omitempty" json:"education,omitempty"`
	ClinicIDs       []primitive.ObjectID `bson:"clinicIds,omitempty" json:"clinicIds,omitempty"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type Education struct {
	Degree      string `bson:"degree" json:"degree"`
	Institution string `bson:"institution" json:"institution"`
	Year        int    `bson:"year" json:"year"`
}

type Booking struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	UserID            primitive.ObjectID  `bson:"userId" json:"userId"`
	DoctorID          primitive.ObjectID  `bson:"doctorId" json:"doctorId"`
	ClinicID          *primitive.ObjectID `bson:"clinicId,omitempty" json:"clinicId,omitempty"`
	SessionType       string              `bson:"sessionType" json:"sessionType"` // Video Call, Audio Call, In-Clinic, Chat
	Date              time.Time           `bson:"date" json:"date"`
	TimeSlot          string              `bson:"timeSlot" json:"timeSlot"`
	Status            string              `bson:"status" json:"status"` // pending, confirmed, completed, cancelled, no-show
	Amount            int                 `bson:"amount" json:"amount"`
	PaymentMethod     string              `bson:"paymentMethod,omitempty" json:"paymentMethod,omitempty"`
	RazorpayOrderID   string              `bson:"razorpayOrderId,omitempty" json:"razorpayOrderId,omitempty"`
	RazorpayPaymentID string              `bson:"razorpayPaymentId,omitempty" json:"razorpayPaymentId,omitempty"`
	PaymentStatus     string              `bson:"paymentStatus" json:"paymentStatus"` // pending, paid, failed, refunded
	Symptoms          []string            `bson:"symptoms,omitempty" json:"symptoms,omitempty"`
	Notes             string              `bson:"notes,omitempty" json:"notes,omitempty"`
	Prescription      *Prescription       `bson:"prescription,omitempty" json:"prescription,omitempty"`
	Diagnosis         []Diagnosis         `bson:"diagnosis,omitempty" json:"diagnosis,omitempty"`
	FollowUpDate      *time.Time          `bson:"followUpDate,omitempty" json:"followUpDate,omitempty"`
	CreatedAt         time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type Prescription struct {
	Medicines    []Medicine `bson:"medicines" json:"medicines"`
	Instructions string     `bson:"instructions,omitempty" json:"instructions,omitempty"`
}

type Medicine struct {
	Name         string `bson:"name" json:"name"`
	Dosage       string `bson:"dosage" json:"dosage"`
	Frequency    string `bson:"frequency" json:"frequency"`
	Duration     string `bson:"duration" json:"duration"`
	Instructions string `bson:"instructions,omitempty" json:"instructions,omitempty"`
}

type Diagnosis struct {
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Date        time.Time `bson:"date" json:"date"`
}

type Journal struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Mood      string             `bson:"mood" json:"mood"`
	Tags      []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	IsPrivate bool               `bson:"isPrivate" json:"isPrivate"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func main() {
	// Load environment variables from either current directory or parent directory
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file from current or parent directory")
		}
	}

	log.Println("Successfully loaded .env file")

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in environment variables")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Get database name from environment or use default
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "hsb_backend"
	}
	db := client.Database(dbName)

	// Drop existing collections to start fresh
	collections := []string{"users", "doctors", "bookings", "journals", "clinics", "appointments", "prescriptions", "medical_records", "pregnancy_trackers", "symptoms", "period_cycles", "mental_health_logs"}
	for _, coll := range collections {
		err = db.Collection(coll).Drop(context.Background())
		if err != nil {
			log.Printf("Warning: Failed to drop collection %s: %v", coll, err)
		} else {
			log.Printf("Dropped collection: %s", coll)
		}
	}

	// Insert test users
	dobUser1, _ := time.Parse("2006-01-02", "1990-05-15")
	dobUser2, _ := time.Parse("2006-01-02", "1988-11-22")
	dobUser3, _ := time.Parse("2006-01-02", "1995-03-10")

	users := []User{
		{
			ID:           objectIDFromHex(TestUser1ID),
			Name:         "Dr. Priya Sharma",
			Email:        "dr.priyasharma@example.com",
			PhoneNumber:  "+919876543210",
			Password:     hashPassword("password123"),
			DateOfBirth:  dobUser1,
			Gender:       "female",
			IsVerified:   true,
			IsActive:     true,
			ProfileImage: "https://example.com/profiles/dr_priya.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           objectIDFromHex(TestUser2ID),
			Name:         "Aarav Patel",
			Email:        "aarav.patel@example.com",
			PhoneNumber:  "+919811223344",
			Password:     hashPassword("password123"),
			DateOfBirth:  dobUser2,
			Gender:       "male",
			IsVerified:   true,
			IsActive:     true,
			ProfileImage: "https://example.com/profiles/aarav.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           objectIDFromHex(TestUser3ID),
			Name:         "Ananya Reddy",
			Email:        "ananya.reddy@example.com",
			PhoneNumber:  "+917788990011",
			Password:     hashPassword("password123"),
			DateOfBirth:  dobUser3,
			Gender:       "female",
			IsVerified:   true,
			IsActive:     true,
			ProfileImage: "https://example.com/profiles/ananya.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	userIDs := make([]primitive.ObjectID, len(users))
	for i, user := range users {
		result, err := db.Collection("users").InsertOne(context.Background(), user)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user.Email, err)
			continue
		}
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			userIDs[i] = id
			log.Printf("Inserted user: %s (%s)", user.Email, id.Hex())
		} else {
			log.Printf("Failed to get inserted ID for user: %s", user.Email)
		}
	}

	// Insert test doctors
	doctors := []Doctor{
		{
			UserID:          objectIDFromHex(TestUser1ID),
			Specialties:     []string{"Gynecology", "Obstetrics", "Infertility"},
			Qualifications:  "MD (Obstetrics & Gynecology), DNB, FNB (Reproductive Medicine)",
			Experience:      12,
			Bio:             "Senior Consultant Gynecologist with extensive experience in high-risk pregnancies and infertility treatments. Special interest in minimally invasive surgeries.",
			IsAvailable:     true,
			ConsultationFee: 1500,
			Languages:       []string{"English", "Hindi", "Marathi", "Gujarati"},
			Education: []Education{
				{Degree: "MBBS", Institution: "Grant Medical College, Mumbai", Year: 2008},
				{Degree: "MD (OBGYN)", Institution: "KEM Hospital, Mumbai", Year: 2012},
				{Degree: "FNB (Reproductive Medicine)", Institution: "IRM, Delhi", Year: 2014},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	doctorIDs := make([]primitive.ObjectID, len(doctors))
	for i, doctor := range doctors {
		result, err := db.Collection("doctors").InsertOne(context.Background(), doctor)
		if err != nil {
			log.Printf("Failed to insert doctor for user %s: %v", doctor.UserID.Hex(), err)
			continue
		}
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			doctorIDs[i] = id
			log.Printf("Inserted doctor with ID: %s", id.Hex())
		}
	}

	// Insert test bookings
	followUpDate := time.Now().AddDate(0, 1, 0) // 1 month from now
	bookings := []Booking{
		{
			UserID:        objectIDFromHex(TestUser2ID),
			DoctorID:      doctorIDs[0],
			SessionType:   "Video Call",
			Date:          time.Now().Add(24 * time.Hour),
			TimeSlot:      "10:00 AM - 11:00 AM",
			Status:        "confirmed",
			Amount:        1500,
			PaymentStatus: "paid",
			PaymentMethod: "online",
			Symptoms:      []string{"Irregular periods", "Acne", "Weight gain"},
			Notes:         "Follow-up consultation for PCOS management",
			Prescription: &Prescription{
				Medicines: []Medicine{
					{
						Name:         "Drospirenone and Ethinyl Estradiol",
						Dosage:       "1 tablet",
						Frequency:    "Once daily",
						Duration:     "21 days",
						Instructions: "Take after dinner",
					},
					{
						Name:         "Metformin",
						Dosage:       "500 mg",
						Frequency:    "Twice daily",
						Duration:     "30 days",
						Instructions: "Take with meals",
					},
				},
				Instructions: "Follow up in 1 month. Get fasting blood sugar and lipid profile tests done before next visit.",
			},
			Diagnosis: []Diagnosis{
				{
					Name:        "Polycystic Ovary Syndrome (PCOS)",
					Description: "PCOS with irregular menses and hyperandrogenism",
					Date:        time.Now(),
				},
			},
			FollowUpDate: &followUpDate,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			UserID:        objectIDFromHex(TestUser3ID),
			DoctorID:      doctorIDs[0],
			SessionType:   "In-Clinic",
			Date:          time.Now().Add(48 * time.Hour),
			TimeSlot:      "11:00 AM - 12:00 PM",
			Status:        "confirmed",
			Amount:        1500,
			PaymentStatus: "pending",
			PaymentMethod: "pay_at_clinic",
			Symptoms:      []string{"Pregnancy confirmation", "Nausea", "Fatigue"},
			Notes:         "First trimester consultation",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for i, booking := range bookings {
		result, err := db.Collection("bookings").InsertOne(context.Background(), booking)
		if err != nil {
			log.Printf("Failed to insert booking %d: %v", i+1, err)
			continue
		}
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			log.Printf("Inserted booking with ID: %s", id.Hex())
		}
	}

	// Insert test journals
	journals := []Journal{
		{
			UserID:    objectIDFromHex(TestUser2ID),
			Title:     "PCOS Journey - Week 1",
			Content:   "Started my PCOS management plan today. Doctor has prescribed medication and recommended lifestyle changes. Feeling hopeful but a bit anxious about the side effects. Need to start with 30 minutes of walking daily and cut down on sugar.",
			Mood:      "hopeful",
			Tags:      []string{"PCOS", "health", "treatment"},
			IsPrivate: false,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			UserID:    objectIDFromHex(TestUser2ID),
			Title:     "Symptom Tracker - Day 3",
			Content:   "Noticing some bloating and mild nausea after starting the medication. Doctor said this is normal in the first week. Keeping track of my symptoms as advised. Sleep quality has been poor for the past two nights.",
			Mood:      "tired",
			Tags:      []string{"symptoms", "medication", "PCOS"},
			IsPrivate: true,
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			UserID:    objectIDFromHex(TestUser3ID),
			Title:     "Exciting News!",
			Content:   "Took a home pregnancy test today and it was positive! Feeling a mix of excitement and nervousness. Booked an appointment with Dr. Sharma to confirm and get started on prenatal care. Haven't told anyone yet, not even my partner - want to do something special to share the news.",
			Mood:      "excited",
			Tags:      []string{"pregnancy", "first-trimester", "excited"},
			IsPrivate: true,
			CreatedAt: time.Now().Add(-2 * 24 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * 24 * time.Hour),
		},
	}

	for i, journal := range journals {
		result, err := db.Collection("journals").InsertOne(context.Background(), journal)
		if err != nil {
			log.Printf("Failed to insert journal entry %d: %v", i+1, err)
			continue
		}
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			log.Printf("Inserted journal entry with ID: %s", id.Hex())
		}
	}

	// Insert additional test data for other schemas
	insertClinicData(db)
	insertPregnancyTrackerData(db)
	insertSymptomData(db)
	insertPeriodCycleData(db)
	insertMentalHealthData(db)

	log.Println("✅ All test data inserted successfully!")
}

// Helper function to generate ObjectID from hex string
func objectIDFromHex(hex string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(hex)
	return objID
}

// Helper function to hash passwords
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

// Insert test data for clinics
func insertClinicData(db *mongo.Database) {
	// Create clinic IDs
	clinic1ID := primitive.NewObjectID()
	clinic2ID := primitive.NewObjectID()

	// Insert clinics
	clinics := []interface{}{
		map[string]interface{}{
			"_id":         clinic1ID,
			"name":        "Healing Touch Women's Clinic",
			"address":     "123 Health Street, Mumbai, Maharashtra 400001",
			"phone":       "+912234567890",
			"email":       "info@healingtouch.com",
			"description": "Comprehensive women's healthcare center specializing in gynecology, obstetrics, and fertility treatments.",
			"location": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{72.8777, 19.0760}, // Mumbai coordinates
			},
			"facilities": []string{"Parking", "Pharmacy", "Lab", "Ultrasound", "ECG"},
			"timings": map[string]string{
				"monday":    "9:00 AM - 8:00 PM",
				"tuesday":   "9:00 AM - 8:00 PM",
				"wednesday": "9:00 AM - 8:00 PM",
				"thursday":  "9:00 AM - 8:00 PM",
				"friday":    "9:00 AM - 8:00 PM",
				"saturday":  "9:00 AM - 2:00 PM",
				"sunday":    "Closed",
			},
			"doctors":   []primitive.ObjectID{objectIDFromHex(TestUser1ID)},
			"isActive":  true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		map[string]interface{}{
			"_id":         clinic2ID,
			"name":        "Mother & Child Care Center",
			"address":     "456 Wellness Avenue, Andheri East, Mumbai 400059",
			"phone":       "+912298765432",
			"email":       "care@motherchild.in",
			"description": "Specialized center for pregnancy care, childbirth, and pediatric services with state-of-the-art facilities.",
			"location": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{72.8762, 19.1175}, // Andheri East coordinates
			},
			"facilities": []string{"Parking", "Pharmacy", "NICU", "Labor & Delivery", "Pediatric Care"},
			"timings": map[string]string{
				"monday":    "8:00 AM - 9:00 PM",
				"tuesday":   "8:00 AM - 9:00 PM",
				"wednesday": "8:00 AM - 9:00 PM",
				"thursday":  "8:00 AM - 9:00 PM",
				"friday":    "8:00 AM - 9:00 PM",
				"saturday":  "9:00 AM - 2:00 PM",
				"sunday":    "Emergency Only",
			},
			"doctors":   []primitive.ObjectID{objectIDFromHex(TestUser1ID)},
			"isActive":  true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}

	_, err := db.Collection("clinics").InsertMany(context.Background(), clinics)
	if err != nil {
		log.Printf("Failed to insert clinics: %v", err)
	} else {
		log.Println("✅ Inserted test clinics")
	}

	// Update doctor with clinic IDs
	_, err = db.Collection("doctors").UpdateOne(
		context.Background(),
		map[string]interface{}{"userId": objectIDFromHex(TestUser1ID)},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"clinicIds": []primitive.ObjectID{clinic1ID, clinic2ID},
				"updatedAt": time.Now(),
			},
		},
	)
	if err != nil {
		log.Printf("Failed to update doctor with clinic IDs: %v", err)
	}
}

// Insert test data for pregnancy tracking
func insertPregnancyTrackerData(db *mongo.Database) {
	pregnancyTrackers := []interface{}{
		map[string]interface{}{
			"userId":          objectIDFromHex(TestUser3ID),
			"startDate":       time.Now().AddDate(0, -1, 0), // 1 month ago
			"dueDate":         time.Now().AddDate(0, 8, 0),  // 8 months from now
			"lastPeriodDate":  time.Now().AddDate(0, -2, 0), // 2 months ago
			"pregnancyWeek":   5,
			"trimester":       1,
			"notes":           "First pregnancy. Experiencing mild nausea in the mornings.",
			"doctorId":        objectIDFromHex(TestUser1ID),
			"nextAppointment": time.Now().AddDate(0, 0, 14), // 2 weeks from now
			"isActive":        true,
			"createdAt":       time.Now(),
			"updatedAt":       time.Now(),
		},
	}

	_, err := db.Collection("pregnancy_trackers").InsertMany(context.Background(), pregnancyTrackers)
	if err != nil {
		log.Printf("Failed to insert pregnancy trackers: %v", err)
	} else {
		log.Println("✅ Inserted test pregnancy trackers")
	}
}

// Insert test data for symptom tracking
func insertSymptomData(db *mongo.Database) {
	symptoms := []interface{}{
		// User 2 (PCOS)
		map[string]interface{}{
			"userId":    objectIDFromHex(TestUser2ID),
			"name":      "Irregular Periods",
			"severity":  "moderate",
			"frequency": "frequent",
			"notes":     "Periods are irregular, occurring every 35-60 days",
			"date":      time.Now().AddDate(0, -1, 0),
			"tags":      []string{"PCOS", "menstrual"},
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		map[string]interface{}{
			"userId":    objectIDFromHex(TestUser2ID),
			"name":      "Acne Breakout",
			"severity":  "mild",
			"frequency": "occasional",
			"location":  "chin and jawline",
			"date":      time.Now().AddDate(0, 0, -5),
			"tags":      []string{"PCOS", "skin"},
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		// User 3 (Pregnancy)
		map[string]interface{}{
			"userId":    objectIDFromHex(TestUser3ID),
			"name":      "Morning Sickness",
			"severity":  "moderate",
			"frequency": "daily",
			"notes":     "Feeling nauseous mostly in the mornings, sometimes throughout the day",
			"date":      time.Now().AddDate(0, 0, -3),
			"tags":      []string{"pregnancy", "first-trimester"},
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		map[string]interface{}{
			"userId":    objectIDFromHex(TestUser3ID),
			"name":      "Fatigue",
			"severity":  "moderate",
			"frequency": "daily",
			"notes":     "Feeling extremely tired, especially in the afternoons",
			"date":      time.Now().AddDate(0, 0, -2),
			"tags":      []string{"pregnancy", "first-trimester"},
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}

	_, err := db.Collection("symptoms").InsertMany(context.Background(), symptoms)
	if err != nil {
		log.Printf("Failed to insert symptoms: %v", err)
	} else {
		log.Println("✅ Inserted test symptoms")
	}
}

// Insert test data for period cycle tracking
func insertPeriodCycleData(db *mongo.Database) {
	periodCycles := []interface{}{
		// User 2 (PCOS)
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser2ID),
			"startDate":    time.Now().AddDate(0, -3, 0), // 3 months ago
			"endDate":      time.Now().AddDate(0, -3, 5), // 5 days later
			"cycleLength":  45,                           // Irregular cycle
			"periodLength": 6,
			"flow":         "medium",
			"symptoms":     []string{"cramps", "bloating"},
			"mood":         []string{"irritable", "tired"},
			"notes":        "Heavy flow on day 2, took pain medication",
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser2ID),
			"startDate":    time.Now().AddDate(0, -1, 10), // 1 month and 10 days ago
			"endDate":      time.Now().AddDate(0, -1, 15), // 5 days later
			"cycleLength":  48,                            // Irregular cycle
			"periodLength": 5,
			"flow":         "light",
			"symptoms":     []string{"bloating", "headache"},
			"mood":         []string{"anxious"},
			"notes":        "Lighter than usual period, followed by spotting",
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		// User 3 (Pregnancy - no recent periods)
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser3ID),
			"startDate":    time.Now().AddDate(0, -3, 0), // 3 months ago
			"endDate":      time.Now().AddDate(0, -3, 5), // 5 days later
			"cycleLength":  28,
			"periodLength": 5,
			"flow":         "medium",
			"symptoms":     []string{"cramps"},
			"mood":         []string{"normal"},
			"notes":        "Regular period before pregnancy",
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
	}

	_, err := db.Collection("period_cycles").InsertMany(context.Background(), periodCycles)
	if err != nil {
		log.Printf("Failed to insert period cycles: %v", err)
	} else {
		log.Println("✅ Inserted test period cycles")
	}
}

// Insert test data for mental health tracking
func insertMentalHealthData(db *mongo.Database) {
	mentalHealthLogs := []interface{}{
		// User 2 (PCOS)
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser2ID),
			"date":         time.Now().AddDate(0, 0, -2),
			"mood":         "anxious",
			"anxietyLevel": 6,
			"stressLevel":  7,
			"sleepQuality": 5,
			"energyLevel":  4,
			"notes":        "Feeling stressed about PCOS diagnosis. Worried about fertility issues.",
			"tags":         []string{"PCOS", "anxiety", "stress"},
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser2ID),
			"date":         time.Now(),
			"mood":         "hopeful",
			"anxietyLevel": 4,
			"stressLevel":  5,
			"sleepQuality": 6,
			"energyLevel":  5,
			"notes":        "Had a good consultation with Dr. Sharma. Feeling more positive about treatment plan.",
			"tags":         []string{"PCOS", "treatment", "doctor-visit"},
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		// User 3 (Pregnancy)
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser3ID),
			"date":         time.Now().AddDate(0, 0, -5),
			"mood":         "excited",
			"anxietyLevel": 3,
			"stressLevel":  4,
			"sleepQuality": 7,
			"energyLevel":  6,
			"notes":        "Found out I'm pregnant! So excited but also nervous about the journey ahead.",
			"tags":         []string{"pregnancy", "excitement", "first-trimester"},
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		map[string]interface{}{
			"userId":       objectIDFromHex(TestUser3ID),
			"date":         time.Now(),
			"mood":         "tired",
			"anxietyLevel": 5,
			"stressLevel":  6,
			"sleepQuality": 4,
			"energyLevel":  3,
			"notes":        "Morning sickness is really taking a toll. Having trouble keeping food down. Need to talk to doctor about solutions.",
			"tags":         []string{"pregnancy", "morning-sickness", "fatigue"},
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
	}

	_, err := db.Collection("mental_health_logs").InsertMany(context.Background(), mentalHealthLogs)
	if err != nil {
		log.Printf("Failed to insert mental health logs: %v", err)
	} else {
		log.Println("✅ Inserted test mental health logs")
	}
}
