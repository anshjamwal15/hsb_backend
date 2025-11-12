package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type PeriodService struct {
	periodRepo repositories.PeriodRepository
}

func NewPeriodService(periodRepo repositories.PeriodRepository) *PeriodService {
	return &PeriodService{
		periodRepo: periodRepo,
	}
}

func (s *PeriodService) AddPeriodCycle(ctx context.Context, cycle *entities.PeriodCycle) error {
	return s.periodRepo.Create(ctx, cycle)
}

func (s *PeriodService) GetPeriodCycles(ctx context.Context, userID string) ([]*entities.PeriodCycle, error) {
	return s.periodRepo.FindByUserID(ctx, userID, 12) // Last 12 cycles
}

func (s *PeriodService) ResetPeriodTracker(ctx context.Context, userID string) error {
	return s.periodRepo.DeleteByUserID(ctx, userID)
}
