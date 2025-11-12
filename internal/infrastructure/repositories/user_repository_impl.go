package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryImpl struct {
	collection    *mongo.Collection
	otpCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repositories.UserRepository {
	return &userRepositoryImpl{
		collection:    db.Collection("users"),
		otpCollection: db.Collection("otps"),
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.IsVerified = false

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id, "isActive": true}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email, "isActive": true}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"isActive": false, "updatedAt": time.Now()}},
	)
	return err
}

func (r *userRepositoryImpl) CreateOTP(ctx context.Context, otp *entities.OTP) error {
	otp.CreatedAt = time.Now()
	_, err := r.otpCollection.InsertOne(ctx, otp)
	return err
}

func (r *userRepositoryImpl) FindOTPByEmail(ctx context.Context, email string) (*entities.OTP, error) {
	var otp entities.OTP
	err := r.otpCollection.FindOne(ctx, bson.M{"email": email}).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *userRepositoryImpl) DeleteOTP(ctx context.Context, email string) error {
	_, err := r.otpCollection.DeleteOne(ctx, bson.M{"email": email})
	return err
}
