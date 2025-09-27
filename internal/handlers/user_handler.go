package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"topup-backend/internal/services"
	"topup-backend/internal/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	services  *services.ServiceContainer
	validator *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(services *services.ServiceContainer) *UserHandler {
	return &UserHandler{
		services:  services,
		validator: validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Registration request"
// @Success 201 {object} utils.APIResponse{data=services.LoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Register user
	response, err := h.services.UserService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login request"
// @Success 200 {object} utils.APIResponse{data=services.LoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Login user
	response, err := h.services.UserService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token request"
// @Success 200 {object} utils.APIResponse{data=services.LoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Refresh token
	response, err := h.services.UserService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client should discard tokens)
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	// In JWT implementation, logout is handled on client side by discarding tokens
	// In production, you might want to maintain a blacklist of tokens
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=models.User}
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /auth/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get user profile
	user, err := h.services.UserService.GetProfile(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} utils.APIResponse{data=models.User}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Update profile
	user, err := h.services.UserService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user password
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.ChangePasswordRequest true "Change password request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		validationErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Change password
	if err := h.services.UserService.ChangePassword(userID.(uint), &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// GetUserByID godoc
// @Summary Get user by ID (Admin only)
// @Description Get user details by ID for admin purposes
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.APIResponse{data=models.User}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get user profile
	user, err := h.services.UserService.GetProfile(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}
