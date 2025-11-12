package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type journalRepositoryImpl struct {
	collection *mongo.Collection
}

func NewJournalRepository(db *mongo.Database) repositories.JournalRepository {
	return &journalRepositoryImpl{
		collection: db.Collection("journals"),
	}
}

func (r *journalRepositoryImpl) Create(ctx context.Context, journal *entities.Journal) error {
	journal.CreatedAt = time.Now()
	journal.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, journal)
	if err != nil {
		return err
	}
	journal.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *journalRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Journal, error) {
	var journal entities.Journal
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&journal)
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func (r *journalRepositoryImpl) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int, search, category string, date *time.Time) ([]*entities.Journal, int64, error) {
	filter := bson.M{"userId": userID}

	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"content": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if category != "" {
		filter["category"] = category
	}

	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		filter["createdAt"] = bson.M{"$gte": startOfDay, "$lt": endOfDay}
	}

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var journals []*entities.Journal
	if err := cursor.All(ctx, &journals); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return journals, total, nil
}

func (r *journalRepositoryImpl) Update(ctx context.Context, journal *entities.Journal) error {
	journal.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": journal.ID},
		bson.M{"$set": journal},
	)
	return err
}

func (r *journalRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
