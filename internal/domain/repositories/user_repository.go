package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	// OTP operations
	CreateOTP(ctx context.Context, otp *entities.OTP) error
	FindOTPByEmail(ctx context.Context, email string) (*entities.OTP, error)
	DeleteOTP(ctx context.Context, email string) error
}
