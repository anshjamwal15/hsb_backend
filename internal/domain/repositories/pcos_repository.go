package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type PCOSRepository interface {
	Create(ctx context.Context, assessment *entities.PCOSAssessment) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.PCOSAssessment, error)
	FindLatestByUserID(ctx context.Context, userID string) (*entities.PCOSAssessment, error)
}
