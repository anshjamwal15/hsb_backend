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

type PCOSRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPCOSRepository(db *mongo.Database) *PCOSRepositoryImpl {
	return &PCOSRepositoryImpl{
		collection: db.Collection("pcos_assessments"),
	}
}

func (r *PCOSRepositoryImpl) Create(ctx context.Context, assessment *entities.PCOSAssessment) error {
	now := time.Now()
	assessment.CreatedAt = now
	assessment.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, assessment)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		assessment.ID = oid
	}

	return nil
}

func (r *PCOSRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.PCOSAssessment, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assessments []*entities.PCOSAssessment
	if err = cursor.All(ctx, &assessments); err != nil {
		return nil, err
	}

	return assessments, nil
}

func (r *PCOSRepositoryImpl) FindLatestByUserID(ctx context.Context, userID string) (*entities.PCOSAssessment, error) {
	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userOID}
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var assessment entities.PCOSAssessment
	err = r.collection.FindOne(ctx, filter, opts).Decode(&assessment)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &assessment, nil
}
