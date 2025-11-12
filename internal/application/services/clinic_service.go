package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type ClinicService struct {
	clinicRepo repositories.ClinicRepository
}

func NewClinicService(clinicRepo repositories.ClinicRepository) *ClinicService {
	return &ClinicService{
		clinicRepo: clinicRepo,
	}
}

func (s *ClinicService) GetClinics(ctx context.Context) ([]*entities.Clinic, error) {
	return s.clinicRepo.FindAll(ctx)
}

func (s *ClinicService) CreateBooking(ctx context.Context, booking *entities.ClinicBooking) error {
	booking.PaymentStatus = "pending"
	booking.Status = "pending"
	return s.clinicRepo.CreateBooking(ctx, booking)
}

func (s *ClinicService) GetMyBookings(ctx context.Context, userID string) ([]*entities.ClinicBooking, error) {
	return s.clinicRepo.FindBookingsByUserID(ctx, userID)
}

func (s *ClinicService) VerifyPayment(ctx context.Context, bookingID, paymentID string) error {
	return s.clinicRepo.UpdateBookingPayment(ctx, bookingID, paymentID, "completed")
}
