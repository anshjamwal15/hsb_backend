package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type WeightService struct {
	weightRepo repositories.WeightRepository
}

func NewWeightService(weightRepo repositories.WeightRepository) *WeightService {
	return &WeightService{
		weightRepo: weightRepo,
	}
}

func (s *WeightService) AddEntry(ctx context.Context, entry *entities.WeightMetabolic) error {
	// Calculate BMI if height is provided
	if entry.Height > 0 && entry.Weight > 0 {
		heightInMeters := entry.Height / 100
		entry.BMI = entry.Weight / (heightInMeters * heightInMeters)
	}

	return s.weightRepo.Create(ctx, entry)
}

func (s *WeightService) GetData(ctx context.Context, userID string) ([]*entities.WeightMetabolic, error) {
	return s.weightRepo.FindByUserID(ctx, userID, 30) // Last 30 entries
}
