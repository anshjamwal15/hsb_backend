package mongodb

import (
    "context"
    "time"

    "github.com/anshjamwal15/hsb_backend/internal/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
    return &userRepository{
        collection: db.Collection("users"),
    }
}

func (r *userRepository) Create(user *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    user.CreatedAt = time.Now().Unix()
    user.UpdatedAt = time.Now().Unix()

    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }

    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        user.ID = oid.Hex()
    }

    return nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user domain.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return nil, domain.ErrDatabase
    }
    return &user, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user domain.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(user.ID)
    if err != nil {
        return err
    }

    user.UpdatedAt = time.Now().Unix()
    update := bson.M{
        "$set": user,
    }

    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    return err
}

func (r *userRepository) Delete(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
    return err
}
