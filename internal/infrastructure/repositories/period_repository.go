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

type PeriodRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPeriodRepository(db *mongo.Database) *PeriodRepositoryImpl {
	return &PeriodRepositoryImpl{
		collection: db.Collection("period_cycles"),
	}
}

func (r *PeriodRepositoryImpl) Create(ctx context.Context, cycle *entities.PeriodCycle) error {
	now := time.Now()
	cycle.CreatedAt = now
	cycle.UpdatedAt = now

	// Ensure UserID is set as ObjectID
	if cycle.UserID.IsZero() {
		userOID, err := primitive.ObjectIDFromHex(cycle.UserID.Hex())
		if err != nil {
			return err
		}
		cycle.UserID = userOID
	}

	result, err := r.collection.InsertOne(ctx, cycle)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		cycle.ID = oid
	}

	return nil
}

func (r *PeriodRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.PeriodCycle, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}
	opts := options.Find().SetSort(bson.D{{Key: "startDate", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cycles []*entities.PeriodCycle
	if err = cursor.All(ctx, &cycles); err != nil {
		return nil, err
	}

	return cycles, nil
}

func (r *PeriodRepositoryImpl) DeleteByUserID(ctx context.Context, userID string) error {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"userId": userOID}
	_, err = r.collection.DeleteMany(ctx, filter)
	return err
}
