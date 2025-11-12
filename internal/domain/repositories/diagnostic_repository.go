package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type DiagnosticRepository interface {
	FindAll(ctx context.Context) ([]*entities.Diagnostic, error)
	CreateBooking(ctx context.Context, booking *entities.DiagnosticBooking) error
	FindBookingsByUserID(ctx context.Context, userID string) ([]*entities.DiagnosticBooking, error)
	UpdateBookingPayment(ctx context.Context, bookingID, paymentID, status string) error
}
