package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type MentalHealthRepository interface {
	CreateResult(ctx context.Context, result *entities.TestResult) error
	FindResultsByUserID(ctx context.Context, userID string, testName string) ([]*entities.TestResult, error)
}
