package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SymptomsRepositoryImpl struct {
	collection *mongo.Collection
}

func NewSymptomsRepository(db *mongo.Database) *SymptomsRepositoryImpl {
	return &SymptomsRepositoryImpl{
		collection: db.Collection("symptoms_tracking"),
	}
}

func (r *SymptomsRepositoryImpl) Create(ctx context.Context, tracking *entities.SymptomsTracking) error {
	now := time.Now()
	tracking.CreatedAt = now
	tracking.UpdatedAt = now

	// Ensure UserID is set as ObjectID
	if tracking.UserID.IsZero() {
		userOID, err := primitive.ObjectIDFromHex(tracking.UserID.Hex())
		if err != nil {
			return err
		}
		tracking.UserID = userOID
	}

	// Ensure Date is set
	if tracking.Date.IsZero() {
		tracking.Date = now
	}

	result, err := r.collection.InsertOne(ctx, tracking)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		tracking.ID = oid
	}

	return nil
}

func (r *SymptomsRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.SymptomsTracking, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trackings []*entities.SymptomsTracking
	if err = cursor.All(ctx, &trackings); err != nil {
		return nil, err
	}

	return trackings, nil
}
