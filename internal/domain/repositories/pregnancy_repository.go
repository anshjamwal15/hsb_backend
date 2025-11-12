package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type PregnancyRepository interface {
	Create(ctx context.Context, tracker *entities.PregnancyTracker) error
	FindByUserID(ctx context.Context, userID string) (*entities.PregnancyTracker, error)
	Update(ctx context.Context, tracker *entities.PregnancyTracker) error
}
