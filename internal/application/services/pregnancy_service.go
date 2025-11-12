package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PregnancyService struct {
	pregnancyRepo repositories.PregnancyRepository
}

func NewPregnancyService(pregnancyRepo repositories.PregnancyRepository) *PregnancyService {
	return &PregnancyService{
		pregnancyRepo: pregnancyRepo,
	}
}

func (s *PregnancyService) GetPregnancyData(ctx context.Context, userID string) (*entities.PregnancyTracker, error) {
	return s.pregnancyRepo.FindByUserID(ctx, userID)
}

func (s *PregnancyService) AddPregnancyEntry(ctx context.Context, tracker *entities.PregnancyTracker) error {
	// Convert userID string to ObjectID if needed
	if tracker.UserID.IsZero() {
		userOID, err := primitive.ObjectIDFromHex(tracker.UserID.Hex())
		if err != nil {
			return err
		}
		tracker.UserID = userOID
	}

	existing, err := s.pregnancyRepo.FindByUserID(ctx, tracker.UserID.Hex())
	if err != nil {
		return err
	}

	if existing != nil {
		tracker.ID = existing.ID
		tracker.CreatedAt = existing.CreatedAt
		return s.pregnancyRepo.Update(ctx, tracker)
	}

	return s.pregnancyRepo.Create(ctx, tracker)
}
