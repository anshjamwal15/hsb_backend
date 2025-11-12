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

type MentalHealthRepositoryImpl struct {
	collection *mongo.Collection
}

func NewMentalHealthRepository(db *mongo.Database) *MentalHealthRepositoryImpl {
	return &MentalHealthRepositoryImpl{
		collection: db.Collection("mental_health_results"),
	}
}

func (r *MentalHealthRepositoryImpl) CreateResult(ctx context.Context, result *entities.TestResult) error {
	now := time.Now()
	result.CreatedAt = now
	result.UpdatedAt = now

	// Set test date if not provided
	if result.TestDate.IsZero() {
		result.TestDate = now
	}

	res, err := r.collection.InsertOne(ctx, result)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		result.ID = oid
	}

	return nil
}

func (r *MentalHealthRepositoryImpl) FindResultsByUserID(ctx context.Context, userID string, testName string) ([]*entities.TestResult, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}
	if testName != "" {
		filter["testName"] = testName
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*entities.TestResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
