package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"github.com/anshjamwal15/hsb_backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, phoneNumber, password string) (*entities.User, string, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &entities.User{
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.Hex(), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.Hex(), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate OTP
	otp := generateOTP()
	otpEntity := &entities.OTP{
		Email:     user.Email,
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	// Delete existing OTP
	s.userRepo.DeleteOTP(ctx, email)

	// Save OTP
	if err := s.userRepo.CreateOTP(ctx, otpEntity); err != nil {
		return err
	}

	// TODO: Send email with OTP
	fmt.Printf("OTP for %s: %s\n", email, otp)

	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, email, code string) error {
	otp, err := s.userRepo.FindOTPByEmail(ctx, email)
	if err != nil {
		return errors.New("invalid OTP")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP expired")
	}

	if otp.Code != code {
		return errors.New("invalid OTP")
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	// Verify OTP
	if err := s.VerifyOTP(ctx, email, otp); err != nil {
		return err
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete OTP
	s.userRepo.DeleteOTP(ctx, email)

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	objID, err := auth.ParseObjectID(userID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, objID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
