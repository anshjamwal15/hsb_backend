package services

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"github.com/anshjamwal15/hsb_backend/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalService struct {
	journalRepo repositories.JournalRepository
}

func NewJournalService(journalRepo repositories.JournalRepository) *JournalService {
	return &JournalService{
		journalRepo: journalRepo,
	}
}

func (s *JournalService) CreateJournal(ctx context.Context, userID, title, content, category string) (*entities.Journal, error) {
	userObjID, err := auth.ParseObjectID(userID)
	if err != nil {
		return nil, err
	}

	journal := &entities.Journal{
		UserID:   userObjID,
		Title:    title,
		Content:  content,
		Category: category,
	}

	if err := s.journalRepo.Create(ctx, journal); err != nil {
		return nil, err
	}

	return journal, nil
}

func (s *JournalService) GetJournals(ctx context.Context, userID string, page, limit int, search, category string, date *time.Time) ([]*entities.Journal, int64, error) {
	userObjID, err := auth.ParseObjectID(userID)
	if err != nil {
		return nil, 0, err
	}

	return s.journalRepo.FindByUserID(ctx, userObjID, page, limit, search, category, date)
}

func (s *JournalService) GetJournalByID(ctx context.Context, journalID string) (*entities.Journal, error) {
	objID, err := primitive.ObjectIDFromHex(journalID)
	if err != nil {
		return nil, err
	}

	return s.journalRepo.FindByID(ctx, objID)
}

func (s *JournalService) UpdateJournal(ctx context.Context, journalID, title, content, category string) (*entities.Journal, error) {
	objID, err := primitive.ObjectIDFromHex(journalID)
	if err != nil {
		return nil, err
	}

	journal, err := s.journalRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		journal.Title = title
	}
	if content != "" {
		journal.Content = content
	}
	if category != "" {
		journal.Category = category
	}

	if err := s.journalRepo.Update(ctx, journal); err != nil {
		return nil, err
	}

	return journal, nil
}

func (s *JournalService) DeleteJournal(ctx context.Context, journalID string) error {
	objID, err := primitive.ObjectIDFromHex(journalID)
	if err != nil {
		return err
	}

	return s.journalRepo.Delete(ctx, objID)
}
