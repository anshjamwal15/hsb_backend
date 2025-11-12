package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type PeriodRepository interface {
	Create(ctx context.Context, cycle *entities.PeriodCycle) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.PeriodCycle, error)
	DeleteByUserID(ctx context.Context, userID string) error
}
