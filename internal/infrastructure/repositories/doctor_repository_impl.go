package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type doctorRepositoryImpl struct {
	db *mongo.Database
}

func NewDoctorRepository(db *mongo.Database) repositories.DoctorRepository {
	return &doctorRepositoryImpl{
		db: db,
	}
}

func (r *doctorRepositoryImpl) collection() *mongo.Collection {
	return r.db.Collection("doctors")
}

func (r *doctorRepositoryImpl) timeSlotCollection() *mongo.Collection {
	return r.db.Collection("doctor_availability")
}

// Create creates a new doctor
func (r *doctorRepositoryImpl) Create(ctx context.Context, doctor *entities.Doctor) (*entities.Doctor, error) {
	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now
	
	// Initialize pointer fields if nil
	if doctor.IsDeleted == nil {
		isDeleted := false
		doctor.IsDeleted = &isDeleted
	}
	if doctor.IsApproved == nil {
		isApproved := false // Default to false, needs admin approval
		doctor.IsApproved = &isApproved
	}
	if doctor.IsAvailable == nil {
		isAvailable := true
		doctor.IsAvailable = &isAvailable
	}

	result, err := r.collection().InsertOne(ctx, doctor)
	if err != nil {
		return nil, fmt.Errorf("failed to create doctor: %w", err)
	}

	// Update the ID of the created doctor
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		doctor.ID = oid.Hex()
	}

	return doctor, nil
}

// FindByID finds a doctor by ID
func (r *doctorRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Doctor, error) {
	var doctor entities.Doctor

	// Convert string ID to ObjectID if needed
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID format: %w", err)
	}

	err = r.collection().FindOne(ctx, bson.M{
		"_id":       objectID,
		"isDeleted": bson.M{"$ne": true},
	}).Decode(&doctor)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find doctor: %w", err)
	}

	return &doctor, nil
}

// FindByEmail finds a doctor by email
func (r *doctorRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.Doctor, error) {
	var doctor entities.Doctor
	err := r.collection().FindOne(ctx, bson.M{
		"email":     email,
		"isDeleted": false,
	}).Decode(&doctor)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find doctor by email: %w", err)
	}

	return &doctor, nil
}

// Update updates a doctor
func (r *doctorRepositoryImpl) Update(ctx context.Context, id string, doctor *entities.Doctor) (*entities.Doctor, error) {
	doctor.UpdatedAt = time.Now()

	update := bson.M{
		"$set": doctor,
	}

	result := r.collection().FindOneAndUpdate(
		ctx,
		bson.M{"_id": id, "isDeleted": false},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("failed to update doctor: %w", err)
	}

	var updatedDoc entities.Doctor
	if err := result.Decode(&updatedDoc); err != nil {
		return nil, fmt.Errorf("failed to decode updated doctor: %w", err)
	}

	return &updatedDoc, nil
}

// Delete marks a doctor as deleted
func (r *doctorRepositoryImpl) Delete(ctx context.Context, id string) error {
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": id, "isDeleted": false},
		update,
	)

	if err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}

	return nil
}

// FindAll finds all doctors with optional filters
func (r *doctorRepositoryImpl) FindAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entities.Doctor, error) {
	// Start with base filter
	filter := bson.M{"isDeleted": false, "isApproved": true}

	// Apply additional filters
	for key, value := range filters {
		filter[key] = value
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find doctors: %w", err)
	}
	defer cursor.Close(ctx)

	var doctors []*entities.Doctor
	if err := cursor.All(ctx, &doctors); err != nil {
		return nil, fmt.Errorf("failed to decode doctors: %w", err)
	}

	return doctors, nil
}

// UpdateAvailability updates a doctor's availability
func (r *doctorRepositoryImpl) UpdateAvailability(ctx context.Context, doctorID string, slots []entities.TimeSlot) error {
	// First, delete existing slots for this doctor
	_, err := r.timeSlotCollection().DeleteMany(ctx, bson.M{"doctor_id": doctorID})
	if err != nil {
		return fmt.Errorf("failed to clear existing availability: %w", err)
	}

	// Insert new slots
	if len(slots) > 0 {
		docs := make([]interface{}, len(slots))
		now := time.Now()

		for i, slot := range slots {
			slot.ID = "" // Clear ID to let MongoDB generate a new one
			slot.DoctorID = doctorID
			slot.CreatedAt = now
			slot.UpdatedAt = now
			docs[i] = slot
		}

		_, err = r.timeSlotCollection().InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to insert availability slots: %w", err)
		}
	}

	return nil
}

// GetAvailability gets a doctor's availability
func (r *doctorRepositoryImpl) GetAvailability(ctx context.Context, doctorID string) ([]entities.TimeSlot, error) {
	cursor, err := r.timeSlotCollection().Find(ctx, bson.M{
		"doctor_id":  doctorID,
		"is_available": true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get availability: %w", err)
	}
	defer cursor.Close(ctx)

	var slots []entities.TimeSlot
	if err := cursor.All(ctx, &slots); err != nil {
		return nil, fmt.Errorf("failed to decode availability: %w", err)
	}

	return slots, nil
}

// FindBySpecialization finds doctors by specialization
func (r *doctorRepositoryImpl) FindBySpecialization(ctx context.Context, specialization string) ([]*entities.Doctor, error) {
	filter := bson.M{
		"specialization": specialization,
		"isDeleted":     false,
		"isApproved":    true,
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find doctors by specialization: %w", err)
	}
	defer cursor.Close(ctx)

	var doctors []*entities.Doctor
	if err := cursor.All(ctx, &doctors); err != nil {
		return nil, fmt.Errorf("failed to decode doctors: %w", err)
	}

	return doctors, nil
}
