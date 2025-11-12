package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type SymptomsService struct {
	symptomsRepo repositories.SymptomsRepository
}

func NewSymptomsService(symptomsRepo repositories.SymptomsRepository) *SymptomsService {
	return &SymptomsService{
		symptomsRepo: symptomsRepo,
	}
}

func (s *SymptomsService) SubmitTracking(ctx context.Context, tracking *entities.SymptomsTracking) error {
	return s.symptomsRepo.Create(ctx, tracking)
}

func (s *SymptomsService) GetTrackingHistory(ctx context.Context, userID string) ([]*entities.SymptomsTracking, error) {
	return s.symptomsRepo.FindByUserID(ctx, userID, 30) // Last 30 entries
}
