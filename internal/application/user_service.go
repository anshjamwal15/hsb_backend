package application

import (
	"errors"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

func (s *UserService) Register(user *domain.User) error {
    // Check if user already exists
    existingUser, err := s.userRepo.FindByEmail(user.Email)
    if err != nil {
        return err
    }
    if existingUser != nil {
        return errors.New("user already exists with this email")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    // Set timestamps
    now := time.Now().Unix()
    user.CreatedAt = now
    user.UpdatedAt = now

    return s.userRepo.Create(user)
}

func (s *UserService) Login(email, password string) (*domain.User, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Don't return the password hash
    user.Password = ""
    return user, nil
}

func (s *UserService) GetProfile(id string) (*domain.User, error) {
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    // Don't return the password hash
    user.Password = ""
    return user, nil
}

func (s *UserService) UpdateProfile(user *domain.User) error {
    // Get existing user to ensure they exist
    existingUser, err := s.userRepo.FindByID(user.ID)
    if err != nil {
        return err
    }
    if existingUser == nil {
        return errors.New("user not found")
    }

    // Don't update password here, use ChangePassword instead
    user.Password = existingUser.Password
    user.UpdatedAt = time.Now().Unix()

    return s.userRepo.Update(user)
}

func (s *UserService) ChangePassword(userID, currentPassword, newPassword string) error {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("user not found")
    }

    // Verify current password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
    if err != nil {
        return errors.New("current password is incorrect")
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)
    user.UpdatedAt = time.Now().Unix()

    return s.userRepo.Update(user)
}
