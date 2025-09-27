package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"topup-backend/internal/config"
	"topup-backend/internal/models"
	"topup-backend/internal/utils"
)

// UserService handles user-related business logic
type UserService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, cfg *config.Config) *UserService {
	return &UserService{
		db:  db,
		cfg: cfg,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
}

// Register registers a new user
func (s *UserService) Register(req *RegisterRequest) (*LoginResponse, error) {
	// Validate phone number
	if !utils.IsValidPhone(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// Normalize phone number
	normalizedPhone := utils.NormalizePhone(req.Phone)

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("phone = ?", normalizedPhone).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this phone number already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Name:         req.Name,
		Phone:        normalizedPhone,
		PasswordHash: hashedPassword,
		Status:       "active",
	}

	if req.Email != "" {
		user.Email = &req.Email
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, "user", "customer", s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Validate phone number
	if !utils.IsValidPhone(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// Normalize phone number
	normalizedPhone := utils.NormalizePhone(req.Phone)

	// Find user
	var user models.User
	if err := s.db.Where("phone = ?", normalizedPhone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid phone number or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid phone number or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, "user", "customer", s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes access token using refresh token
func (s *UserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Parse refresh token
	userID, err := utils.ParseRefreshToken(refreshToken, s.cfg)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateJWT(user.ID, "user", "customer", s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetProfile gets user profile by ID
func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &user, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.User, error) {
	// Find user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update fields
	user.Name = req.Name
	if req.Email != "" {
		user.Email = &req.Email
	} else {
		user.Email = nil
	}

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &user, nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	// Find user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check current password
	if !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
