package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type FSFIService struct {
	mentalHealthRepo repositories.MentalHealthRepository
}

func NewFSFIService(mentalHealthRepo repositories.MentalHealthRepository) *FSFIService {
	return &FSFIService{
		mentalHealthRepo: mentalHealthRepo,
	}
}

func (s *FSFIService) GetTest() entities.MentalHealthTest {
	// Common scale options for FSFI questions
	scaleOptions1to5 := []entities.QuestionOption{
		{Value: 1, Label: "1"},
		{Value: 2, Label: "2"},
		{Value: 3, Label: "3"},
		{Value: 4, Label: "4"},
		{Value: 5, Label: "5"},
	}

	scaleOptions0to5 := append([]entities.QuestionOption{
		{Value: 0, Label: "0"},
	}, scaleOptions1to5...)

	return entities.MentalHealthTest{
		Name:        "fsfi",
		DisplayName: "Female Sexual Function Index (FSFI)",
		Description: "Assessment of female sexual function",
		Questions: []entities.TestQuestion{
			{
				QuestionText:  "How often did you feel sexual desire or interest?",
				ResponseType: "scale",
				Options:      scaleOptions1to5,
			},
			{
				QuestionText:  "How would you rate your level of sexual desire or interest?",
				ResponseType: "scale",
				Options:      scaleOptions1to5,
			},
			{
				QuestionText:  "How often did you feel sexually aroused during sexual activity?",
				ResponseType: "scale",
				Options:      scaleOptions0to5,
			},
			{
				QuestionText:  "How would you rate your level of sexual arousal?",
				ResponseType: "scale",
				Options:      scaleOptions0to5,
			},
		},
	}
}

func (s *FSFIService) SubmitTest(ctx context.Context, result *entities.TestResult) error {
	result.TestName = "fsfi"

	// Calculate score
	score := 0
	for _, answer := range result.Answers {
		if val, ok := answer.(float64); ok {
			score += int(val)
		}
	}
	result.Score = score

	// Determine level (FSFI score < 26.55 indicates sexual dysfunction)
	if score >= 27 {
		result.Level = "normal"
	} else {
		result.Level = "dysfunction"
	}

	return s.mentalHealthRepo.CreateResult(ctx, result)
}

func (s *FSFIService) GetMyResults(ctx context.Context, userID string) ([]*entities.TestResult, error) {
	return s.mentalHealthRepo.FindResultsByUserID(ctx, userID, "fsfi")
}
