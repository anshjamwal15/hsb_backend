package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RazorpayClient interface {
	CreateOrder(amount int, currency, receipt string) (string, error)
}

type BookingService struct {
	bookingRepo   repositories.BookingRepository
	doctorRepo    repositories.DoctorRepository
	razorpayClient RazorpayClient
}

func NewBookingService(
	bookingRepo repositories.BookingRepository, 
	doctorRepo repositories.DoctorRepository,
	razorpayClient RazorpayClient,
) *BookingService {
	return &BookingService{
		bookingRepo:    bookingRepo,
		doctorRepo:     doctorRepo,
		razorpayClient: razorpayClient,
	}
}

// CreateBooking creates a new booking
func (s *BookingService) CreateBooking(ctx context.Context, userID, doctorID, sessionType string, date time.Time, timeSlot, notes string) (*entities.Booking, error) {
	// Input validation
	if userID == "" || doctorID == "" {
		return nil, errors.New("user ID and doctor ID are required")
	}

	// Convert string IDs to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	doctorObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID format: %w", err)
	}

	// Validate date (must be in the future)
	if date.Before(time.Now()) {
		return nil, errors.New("booking date must be in the future")
	}

	// Get doctor details
	doctor, err := s.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("error fetching doctor: %w", err)
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	// Check if doctor is available
	if doctor.IsAvailable != nil && !*doctor.IsAvailable {
		return nil, errors.New("doctor is not available for bookings")
	}

	// Calculate amount based on session type
	var amount int
	switch sessionType {
	case "Video Call":
		if doctor.ConsultationFees == nil {
			return nil, errors.New("video call not available with this doctor")
		}
		amount = doctor.ConsultationFees.VideoCall
	case "Audio Call":
		if doctor.ConsultationFees == nil {
			return nil, errors.New("audio call not available with this doctor")
		}
		amount = doctor.ConsultationFees.AudioCall
	case "In-Clinic":
		if doctor.ConsultationFees == nil {
			return nil, errors.New("in-clinic consultation not available with this doctor")
		}
		amount = doctor.ConsultationFees.InClinic
	default:
		return nil, fmt.Errorf("invalid session type: %s", sessionType)
	}

	// Create booking object
	now := time.Now()
	booking := &entities.Booking{
		ID:            primitive.NewObjectID(),
		UserID:        userObjID,
		DoctorID:      doctorObjID,
		SessionType:   sessionType,
		Date:          date,
		TimeSlot:      timeSlot,
		Status:        "pending",
		Amount:        amount,
		PaymentStatus: "pending",
		Notes:         notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Create Razorpay order
	razorpayOrderID, err := s.razorpayClient.CreateOrder(
		booking.Amount * 100, // Convert to paise
		"INR",
		"order_"+booking.ID.Hex(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Razorpay order: %w", err)
	}

	// Update booking with Razorpay order ID
	booking.RazorpayOrderID = razorpayOrderID

	// Save the booking to the database
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return booking, nil
}

// GetUserBookings retrieves a paginated list of bookings for a user
func (s *BookingService) GetUserBookings(ctx context.Context, userID string, page, limit int) ([]*entities.Booking, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID is required")
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.bookingRepo.FindByUserID(ctx, userID, page, limit)
}

// GetActiveBookings retrieves all active bookings for a user
func (s *BookingService) GetActiveBookings(ctx context.Context, userID string) ([]*entities.Booking, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.bookingRepo.FindActiveByUserID(ctx, userID)
}

// GetDoctorBookingsByDate retrieves all bookings for a doctor on a specific date
func (s *BookingService) GetDoctorBookingsByDate(ctx context.Context, doctorID string, date time.Time) ([]*entities.Booking, error) {
	if doctorID == "" {
		return nil, errors.New("doctor ID is required")
	}

	// Set start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// Get all bookings for this doctor on this date
	bookings, err := s.bookingRepo.FindByDoctorID(ctx, doctorID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// VerifyPayment verifies a payment with Razorpay and updates the booking status
func (s *BookingService) VerifyPayment(ctx context.Context, bookingID, razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {
	// Input validation
	if bookingID == "" || razorpayOrderID == "" || razorpayPaymentID == "" || razorpaySignature == "" {
		return errors.New("all payment verification parameters are required")
	}

	// Find the booking
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("error finding booking: %w", err)
	}

	if booking == nil {
		return errors.New("booking not found")
	}

	// Verify the payment with Razorpay
	// Note: In a real implementation, you would verify the signature using Razorpay's secret
	// This is a simplified example
	if razorpaySignature == "" {
		// Update booking status to failed
		err = s.bookingRepo.UpdatePaymentStatus(ctx, bookingID, "failed", razorpayOrderID, razorpayPaymentID)
		if err != nil {
			return fmt.Errorf("payment verification failed: invalid signature, failed to update booking status: %w", err)
		}
		return errors.New("invalid payment signature")
	}

	// Update booking status to paid
	err = s.bookingRepo.UpdatePaymentStatus(ctx, bookingID, "paid", razorpayOrderID, razorpayPaymentID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	// TODO: Send confirmation email to user and doctor
	// TODO: Create calendar event for the appointment

	return nil
}
