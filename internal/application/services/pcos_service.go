package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type PCOSService struct {
	pcosRepo repositories.PCOSRepository
}

func NewPCOSService(pcosRepo repositories.PCOSRepository) *PCOSService {
	return &PCOSService{
		pcosRepo: pcosRepo,
	}
}

func (s *PCOSService) GetQuestions() []entities.PCOSQuestion {
	return []entities.PCOSQuestion{
		{ID: "q1", Question: "Do you have irregular periods?", Type: "yes_no"},
		{ID: "q2", Question: "Do you experience excessive hair growth?", Type: "yes_no"},
		{ID: "q3", Question: "Do you have acne or oily skin?", Type: "yes_no"},
		{ID: "q4", Question: "Have you experienced weight gain?", Type: "yes_no"},
		{ID: "q5", Question: "Do you have difficulty losing weight?", Type: "yes_no"},
		{ID: "q6", Question: "Do you experience hair thinning or hair loss?", Type: "yes_no"},
		{ID: "q7", Question: "Do you have darkening of skin in body folds?", Type: "yes_no"},
		{ID: "q8", Question: "Have you been diagnosed with insulin resistance?", Type: "yes_no"},
	}
}

func (s *PCOSService) SubmitAssessment(ctx context.Context, assessment *entities.PCOSAssessment) error {
	// Calculate score based on responses
	score := 0
	for _, response := range assessment.Responses {
		if answer, ok := response.(bool); ok && answer {
			score++
		} else if answer, ok := response.(string); ok && answer == "yes" {
			score++
		}
	}
	assessment.Score = score

	// Determine result based on score
	if score <= 2 {
		assessment.Result = "low risk"
	} else if score <= 5 {
		assessment.Result = "moderate risk"
	} else {
		assessment.Result = "high risk"
	}

	return s.pcosRepo.Create(ctx, assessment)
}

func (s *PCOSService) GetHistory(ctx context.Context, userID string) ([]*entities.PCOSAssessment, error) {
	return s.pcosRepo.FindByUserID(ctx, userID, 10)
}

func (s *PCOSService) GetLatestAssessment(ctx context.Context, userID string) (*entities.PCOSAssessment, error) {
	return s.pcosRepo.FindLatestByUserID(ctx, userID)
}
