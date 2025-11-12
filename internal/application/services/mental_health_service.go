package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type MentalHealthService struct {
	mentalHealthRepo repositories.MentalHealthRepository
}

func NewMentalHealthService(mentalHealthRepo repositories.MentalHealthRepository) *MentalHealthService {
	return &MentalHealthService{
		mentalHealthRepo: mentalHealthRepo,
	}
}

func (s *MentalHealthService) GetTests() []entities.MentalHealthTest {
	scaleOptions := []entities.QuestionOption{
		{Value: 0, Label: "0"},
		{Value: 1, Label: "1"},
		{Value: 2, Label: "2"},
		{Value: 3, Label: "3"},
	}

	return []entities.MentalHealthTest{
		{
			Name:        "phq9",
			DisplayName: "PHQ-9 Depression Test",
			Description: "Patient Health Questionnaire for depression screening",
			Questions: []entities.TestQuestion{
				{
					QuestionText: "Little interest or pleasure in doing things",
					ResponseType: "scale",
					Options:     scaleOptions,
				},
				{
					QuestionText: "Feeling down, depressed, or hopeless",
					ResponseType: "scale",
					Options:     scaleOptions,
				},
			},
		},
		{
			Name:        "gad7",
			DisplayName: "GAD-7 Anxiety Test",
			Description: "Generalized Anxiety Disorder 7-item scale",
			Questions: []entities.TestQuestion{
				{
					QuestionText: "Feeling nervous, anxious, or on edge",
					ResponseType: "scale",
					Options:     scaleOptions,
				},
				{
					QuestionText: "Not being able to stop or control worrying",
					ResponseType: "scale",
					Options:     scaleOptions,
				},
			},
		},
	}
}

func (s *MentalHealthService) GetTestByName(testName string) *entities.MentalHealthTest {
	tests := s.GetTests()
	for _, test := range tests {
		if test.Name == testName {
			return &test
		}
	}
	return nil
}

func (s *MentalHealthService) SubmitTestResults(ctx context.Context, result *entities.TestResult) error {
	// Calculate score
	score := 0
	for _, answer := range result.Answers {
		if val, ok := answer.(float64); ok {
			score += int(val)
		}
	}
	result.Score = score

	// Determine level based on score
	if score <= 4 {
		result.Level = "minimal"
	} else if score <= 9 {
		result.Level = "mild"
	} else if score <= 14 {
		result.Level = "moderate"
	} else {
		result.Level = "severe"
	}

	return s.mentalHealthRepo.CreateResult(ctx, result)
}

func (s *MentalHealthService) GetTestResults(ctx context.Context, userID string, testName string) ([]*entities.TestResult, error) {
	return s.mentalHealthRepo.FindResultsByUserID(ctx, userID, testName)
}
