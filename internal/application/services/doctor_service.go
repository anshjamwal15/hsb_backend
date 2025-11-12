package services

import (
	"context"
	"errors"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
)

type DoctorService struct {
	doctorRepo repositories.DoctorRepository
}

func NewDoctorService(doctorRepo repositories.DoctorRepository) *DoctorService {
	return &DoctorService{
		doctorRepo: doctorRepo,
	}
}

// CreateDoctor creates a new doctor
func (s *DoctorService) CreateDoctor(ctx context.Context, doctor *entities.Doctor) (*entities.Doctor, error) {
	// Basic validation
	if doctor.Name == "" || doctor.Email == "" || doctor.Specialization == "" {
		return nil, errors.New("name, email, and specialization are required")
	}

	// Check if doctor with email already exists
	existing, err := s.doctorRepo.FindByEmail(ctx, doctor.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("doctor with this email already exists")
	}

	// Set default values
	isApproved := false // Default to false, needs admin approval
	isDeleted := false
	isAvailable := true

	doctor.IsApproved = &isApproved
	doctor.IsDeleted = &isDeleted
	doctor.IsAvailable = &isAvailable

	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now

	// Create doctor
	return s.doctorRepo.Create(ctx, doctor)
}

// GetDoctorByID retrieves a doctor by ID
func (s *DoctorService) GetDoctorByID(ctx context.Context, id string) (*entities.Doctor, error) {
	if id == "" {
		return nil, errors.New("doctor ID is required")
	}
	return s.doctorRepo.FindByID(ctx, id)
}

// UpdateDoctor updates an existing doctor
func (s *DoctorService) UpdateDoctor(ctx context.Context, id string, update *entities.Doctor) (*entities.Doctor, error) {
	if id == "" {
		return nil, errors.New("doctor ID is required")
	}

	// Get existing doctor to ensure they exist
	existing, err := s.doctorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("doctor not found")
	}

	// Update fields that are allowed to be updated
	existing.Name = update.Name
	existing.Email = update.Email
	existing.Phone = update.Phone
	existing.Specialization = update.Specialization
	existing.Bio = update.Bio
	existing.About = update.About
	existing.Location = update.Location
	existing.Languages = update.Languages
	existing.ConsultationFees = update.ConsultationFees
	existing.Timing = update.Timing
	existing.PackageIncludes = update.PackageIncludes
	existing.UpdatedAt = time.Now()

	// Update approval status if provided (admin operation)
	if update.IsApproved != nil {
		existing.IsApproved = update.IsApproved
	}

	return s.doctorRepo.Update(ctx, id, existing)
}

// DeleteDoctor marks a doctor as deleted
func (s *DoctorService) DeleteDoctor(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("doctor ID is required")
	}

	// Verify doctor exists
	existing, err := s.doctorRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("doctor not found")
	}

	return s.doctorRepo.Delete(ctx, id)
}

// ListDoctors retrieves a list of doctors with optional filters
func (s *DoctorService) ListDoctors(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entities.Doctor, error) {
	return s.doctorRepo.FindAll(ctx, filters, limit, offset)
}

// SetAvailability updates a doctor's availability
func (s *DoctorService) SetAvailability(ctx context.Context, doctorID string, slots []entities.TimeSlot) error {
	if doctorID == "" {
		return errors.New("doctor ID is required")
	}

	// Validate time slots
	for _, slot := range slots {
		if slot.StartTime.After(slot.EndTime) || slot.StartTime.Equal(slot.EndTime) {
			return errors.New("invalid time slot: start time must be before end time")
		}
	}

	return s.doctorRepo.UpdateAvailability(ctx, doctorID, slots)
}

// GetAvailableSlots retrieves available time slots for a doctor
func (s *DoctorService) GetAvailableSlots(ctx context.Context, doctorID string, date time.Time) ([]entities.TimeSlot, error) {
	if doctorID == "" {
		return nil, errors.New("doctor ID is required")
	}

	// Get doctor's availability
	slots, err := s.doctorRepo.GetAvailability(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// Filter slots for the requested date
	var availableSlots []entities.TimeSlot
	for _, slot := range slots {
		if slot.StartTime.Year() == date.Year() && 
		   slot.StartTime.Month() == date.Month() && 
		   slot.StartTime.Day() == date.Day() {
			availableSlots = append(availableSlots, slot)
		}
	}

	return availableSlots, nil
}

// GetDoctorsBySpecialization retrieves doctors by their specialization
func (s *DoctorService) GetDoctorsBySpecialization(ctx context.Context, specialization string) ([]*entities.Doctor, error) {
	if specialization == "" {
		return nil, errors.New("specialization is required")
	}

	return s.doctorRepo.FindBySpecialization(ctx, specialization)
}

// GetDoctors retrieves a paginated list of doctors with optional search
func (s *DoctorService) GetDoctors(ctx context.Context, page, limit int, search string) ([]*entities.Doctor, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build filters
	filters := make(map[string]interface{})
	
	// Only show approved and non-deleted doctors
	isApproved := true
	isDeleted := false
	filters["isApproved"] = &isApproved
	filters["isDeleted"] = &isDeleted

	// Add search filter if provided
	if search != "" {
		filters["search"] = search
	}

	// Get doctors
	doctors, err := s.doctorRepo.FindAll(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	// Note: This is a simplified approach. In production, you'd want a separate Count method
	total := int64(len(doctors))
	if len(doctors) == limit {
		// If we got a full page, there might be more
		// This is an approximation - ideally you'd have a Count method in the repository
		total = int64((page + 1) * limit)
	}

	return doctors, total, nil
}
