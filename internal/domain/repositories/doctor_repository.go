package repositories

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
)

type DoctorRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, doctor *entities.Doctor) (*entities.Doctor, error)
	FindByID(ctx context.Context, id string) (*entities.Doctor, error)
	FindByEmail(ctx context.Context, email string) (*entities.Doctor, error)
	Update(ctx context.Context, id string, doctor *entities.Doctor) (*entities.Doctor, error)
	Delete(ctx context.Context, id string) error

	// Query operations
	FindAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entities.Doctor, error)
	
	// Availability operations
	UpdateAvailability(ctx context.Context, doctorID string, slots []entities.TimeSlot) error
	GetAvailability(ctx context.Context, doctorID string) ([]entities.TimeSlot, error)
	
	// Specialization operations
	FindBySpecialization(ctx context.Context, specialization string) ([]*entities.Doctor, error)
}
