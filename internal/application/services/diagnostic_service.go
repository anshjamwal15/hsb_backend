package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type DiagnosticService struct {
	diagnosticRepo repositories.DiagnosticRepository
}

func NewDiagnosticService(diagnosticRepo repositories.DiagnosticRepository) *DiagnosticService {
	return &DiagnosticService{
		diagnosticRepo: diagnosticRepo,
	}
}

func (s *DiagnosticService) GetDiagnostics(ctx context.Context) ([]*entities.Diagnostic, error) {
	return s.diagnosticRepo.FindAll(ctx)
}

func (s *DiagnosticService) CreateBooking(ctx context.Context, booking *entities.DiagnosticBooking) error {
	booking.PaymentStatus = "pending"
	booking.Status = "pending"
	return s.diagnosticRepo.CreateBooking(ctx, booking)
}

func (s *DiagnosticService) GetDiagnosticsUsers(ctx context.Context, userID string) ([]*entities.DiagnosticBooking, error) {
	return s.diagnosticRepo.FindBookingsByUserID(ctx, userID)
}

func (s *DiagnosticService) VerifyPayment(ctx context.Context, bookingID, paymentID string) error {
	return s.diagnosticRepo.UpdateBookingPayment(ctx, bookingID, paymentID, "completed")
}
