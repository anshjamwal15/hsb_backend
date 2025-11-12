package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type WeightRepository interface {
	Create(ctx context.Context, entry *entities.WeightMetabolic) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.WeightMetabolic, error)
}
