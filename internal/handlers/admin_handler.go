package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"topup-backend/internal/services"
	"topup-backend/internal/utils"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	services  *services.ServiceContainer
	validator *validator.Validate
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(services *services.ServiceContainer) *AdminHandler {
	return &AdminHandler{
		services:  services,
		validator: validator.New(),
	}
}

// Login godoc
// @Summary Admin login
// @Description Authenticate admin user and return tokens
// @Tags admin
// @Accept json
// @Produce json
// @Param request body services.AdminLoginRequest true "Admin login request"
// @Success 200 {object} utils.APIResponse{data=services.AdminLoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /admin/auth/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req services.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	response, err := h.services.AdminService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh admin access token
// @Description Refresh admin access token using refresh token
// @Tags admin
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token request"
// @Success 200 {object} utils.APIResponse{data=services.AdminLoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /admin/auth/refresh-token [post]
func (h *AdminHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.services.AdminService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// Logout godoc
// @Summary Admin logout
// @Description Logout admin user (client should discard tokens)
// @Tags admin
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Router /admin/auth/logout [post]
func (h *AdminHandler) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}

// GetProfile godoc
// @Summary Get admin profile
// @Description Get current admin user profile
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=models.AdminUser}
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /admin/auth/profile [get]
func (h *AdminHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Admin not authenticated")
		return
	}

	admin, err := h.services.AdminService.GetProfile(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", admin)
}

// GetAllAdmins godoc
// @Summary Get all admin users
// @Description Get paginated list of admin users (Super Admin/Admin only)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.APIResponse{data=[]models.AdminUser,pagination=utils.Pagination}
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/users [get]
func (h *AdminHandler) GetAllAdmins(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	admins, total, err := h.services.AdminService.GetAllAdmins(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, limit, int(total))
	utils.SuccessResponseWithPagination(c, http.StatusOK, "Admins retrieved successfully", admins, pagination)
}

// CreateAdmin godoc
// @Summary Create admin user
// @Description Create a new admin user (Super Admin/Admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.CreateAdminRequest true "Admin creation request"
// @Success 201 {object} utils.APIResponse{data=models.AdminUser}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/users [post]
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	creatorID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Admin not authenticated")
		return
	}

	var req services.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	admin, err := h.services.AdminService.CreateAdmin(&req, creatorID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin created successfully", admin)
}

// GetAdminByID godoc
// @Summary Get admin by ID
// @Description Get admin user details by ID (Super Admin/Admin only)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} utils.APIResponse{data=models.AdminUser}
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetAdminByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// UpdateAdmin godoc
// @Summary Update admin user
// @Description Update admin user by ID (Super Admin/Admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Param request body services.UpdateAdminRequest true "Admin update request"
// @Success 200 {object} utils.APIResponse{data=models.AdminUser}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/users/{id} [put]
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// DeleteAdmin godoc
// @Summary Delete admin user
// @Description Delete admin user by ID (Super Admin only)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}
