package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PregnancyRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPregnancyRepository(db *mongo.Database) *PregnancyRepositoryImpl {
	return &PregnancyRepositoryImpl{
		collection: db.Collection("pregnancy_trackers"),
	}
}

func (r *PregnancyRepositoryImpl) Create(ctx context.Context, tracker *entities.PregnancyTracker) error {
	now := time.Now()
	tracker.CreatedAt = now
	tracker.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, tracker)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		tracker.ID = oid
	}

	return nil
}

func (r *PregnancyRepositoryImpl) FindByUserID(ctx context.Context, userID string) (*entities.PregnancyTracker, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}

	var tracker entities.PregnancyTracker
	err = r.collection.FindOne(ctx, filter).Decode(&tracker)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tracker, nil
}

func (r *PregnancyRepositoryImpl) Update(ctx context.Context, tracker *entities.PregnancyTracker) error {
	tracker.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"dueDate":        tracker.DueDate,
			"lastPeriodDate": tracker.LastPeriodDate,
			"weight":         tracker.Weight,
			"notes":          tracker.Notes,
			"updatedAt":      tracker.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateByID(ctx, tracker.ID, update)
	return err
}
