package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type BookingRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, booking *entities.Booking) error
	FindByID(ctx context.Context, id string) (*entities.Booking, error)
	Update(ctx context.Context, booking *entities.Booking) error

	// Query operations
	FindByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Booking, int64, error)
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Booking, error)
	
	// New methods for time slot and doctor-based queries
	FindByDoctorAndTimeSlot(ctx context.Context, doctorID string, date time.Time, timeSlot string) ([]*entities.Booking, error)
	FindByDoctorID(ctx context.Context, doctorID string, startDate, endDate time.Time) ([]*entities.Booking, error)
	
	// Payment related
	UpdatePaymentStatus(ctx context.Context, bookingID string, status string, razorpayOrderID, razorpayPaymentID string) error
}
