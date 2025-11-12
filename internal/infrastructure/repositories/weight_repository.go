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

type WeightRepositoryImpl struct {
	collection *mongo.Collection
}

func NewWeightRepository(db *mongo.Database) *WeightRepositoryImpl {
	return &WeightRepositoryImpl{
		collection: db.Collection("weight_metabolic"),
	}
}

func (r *WeightRepositoryImpl) Create(ctx context.Context, entry *entities.WeightMetabolic) error {
	now := time.Now().Unix()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, entry)
	if err != nil {
		return err
	}

	entry.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *WeightRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.WeightMetabolic, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*entities.WeightMetabolic
	if err = cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
