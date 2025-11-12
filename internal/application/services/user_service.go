package services

import (
	"context"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"github.com/anshjamwal15/hsb_backend/pkg/auth"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*entities.User, error) {
	objID, err := auth.ParseObjectID(userID)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, objID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID, name, phoneNumber, profileImage string) (*entities.User, error) {
	objID, err := auth.ParseObjectID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if phoneNumber != "" {
		user.PhoneNumber = phoneNumber
	}
	if profileImage != "" {
		user.ProfileImage = profileImage
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
