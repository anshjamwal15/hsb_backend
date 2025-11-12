package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type SymptomsRepository interface {
	Create(ctx context.Context, tracking *entities.SymptomsTracking) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.SymptomsTracking, error)
}
