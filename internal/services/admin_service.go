package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"topup-backend/internal/config"
	"topup-backend/internal/models"
	"topup-backend/internal/utils"
)

// AdminService handles admin-related business logic
type AdminService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAdminService creates a new admin service
func NewAdminService(db *gorm.DB, cfg *config.Config) *AdminService {
	return &AdminService{
		db:  db,
		cfg: cfg,
	}
}

// AdminLoginRequest represents admin login request
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AdminLoginResponse represents admin login response
type AdminLoginResponse struct {
	Admin        *models.AdminUser `json:"admin"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
}

// CreateAdminRequest represents create admin request
type CreateAdminRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin operator moderator"`
}

// UpdateAdminRequest represents update admin request
type UpdateAdminRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Email  string `json:"email,omitempty" validate:"omitempty,email"`
	Role   string `json:"role" validate:"required,oneof=super_admin admin operator moderator"`
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

// Login authenticates an admin user
func (s *AdminService) Login(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	// Find admin user
	var admin models.AdminUser
	if err := s.db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	// Check if admin is active
	if admin.Status != "active" {
		return nil, errors.New("admin account is not active")
	}

	// Check password
	if !utils.CheckPassword(req.Password, admin.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// Update last login
	now := time.Now()
	admin.LastLoginAt = &now
	s.db.Save(&admin)

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateJWT(admin.ID, "admin", admin.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AdminLoginResponse{
		Admin:        &admin,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes admin access token
func (s *AdminService) RefreshToken(refreshToken string) (*AdminLoginResponse, error) {
	// Parse refresh token
	userID, err := utils.ParseRefreshToken(refreshToken, s.cfg)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find admin
	var admin models.AdminUser
	if err := s.db.First(&admin, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	// Check if admin is active
	if admin.Status != "active" {
		return nil, errors.New("admin account is not active")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateJWT(admin.ID, "admin", admin.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AdminLoginResponse{
		Admin:        &admin,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetProfile gets admin profile by ID
func (s *AdminService) GetProfile(adminID uint) (*models.AdminUser, error) {
	var admin models.AdminUser
	if err := s.db.Preload("Creator").First(&admin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin profile: %w", err)
	}

	return &admin, nil
}

// GetAllAdmins gets all admin users with pagination
func (s *AdminService) GetAllAdmins(page, limit int) ([]models.AdminUser, int64, error) {
	var admins []models.AdminUser
	var total int64

	offset := utils.CalculateOffset(page, limit)

	// Count total
	if err := s.db.Model(&models.AdminUser{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count admins: %w", err)
	}

	// Get admins with pagination
	if err := s.db.Preload("Creator").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&admins).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get admins: %w", err)
	}

	return admins, total, nil
}

// CreateAdmin creates a new admin user
func (s *AdminService) CreateAdmin(req *CreateAdminRequest, creatorID uint) (*models.AdminUser, error) {
	// Check if username already exists
	var existingAdmin models.AdminUser
	if err := s.db.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin
	admin := models.AdminUser{
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		Status:       "active",
		CreatedBy:    &creatorID,
	}

	if req.Email != "" {
		admin.Email = &req.Email
	}

	if err := s.db.Create(&admin).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	// Load creator relationship
	if err := s.db.Preload("Creator").First(&admin, admin.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load admin with creator: %w", err)
	}

	return &admin, nil
}

// UpdateAdmin updates an admin user
func (s *AdminService) UpdateAdmin(adminID uint, req *UpdateAdminRequest) (*models.AdminUser, error) {
	// Find admin
	var admin models.AdminUser
	if err := s.db.First(&admin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	// Update fields
	admin.Name = req.Name
	admin.Role = req.Role
	admin.Status = req.Status

	if req.Email != "" {
		admin.Email = &req.Email
	} else {
		admin.Email = nil
	}

	// Save changes
	if err := s.db.Save(&admin).Error; err != nil {
		return nil, fmt.Errorf("failed to update admin: %w", err)
	}

	// Load relationships
	if err := s.db.Preload("Creator").First(&admin, admin.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load admin with creator: %w", err)
	}

	return &admin, nil
}

// DeleteAdmin deletes an admin user
func (s *AdminService) DeleteAdmin(adminID uint) error {
	// Check if admin exists
	var admin models.AdminUser
	if err := s.db.First(&admin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("admin not found")
		}
		return fmt.Errorf("failed to find admin: %w", err)
	}

	// Cannot delete super_admin
	if admin.Role == "super_admin" {
		return errors.New("cannot delete super admin")
	}

	// Soft delete
	if err := s.db.Delete(&admin).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	return nil
}

// GetAdminByID gets admin by ID
func (s *AdminService) GetAdminByID(adminID uint) (*models.AdminUser, error) {
	var admin models.AdminUser
	if err := s.db.Preload("Creator").First(&admin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return &admin, nil
}
