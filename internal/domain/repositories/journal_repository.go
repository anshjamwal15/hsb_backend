package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalRepository interface {
	Create(ctx context.Context, journal *entities.Journal) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Journal, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int, search, category string, date *time.Time) ([]*entities.Journal, int64, error)
	Update(ctx context.Context, journal *entities.Journal) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
