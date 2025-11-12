package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type ClinicRepository interface {
	FindAll(ctx context.Context) ([]*entities.Clinic, error)
	CreateBooking(ctx context.Context, booking *entities.ClinicBooking) error
	FindBookingsByUserID(ctx context.Context, userID string) ([]*entities.ClinicBooking, error)
	UpdateBookingPayment(ctx context.Context, bookingID, paymentID, status string) error
}
