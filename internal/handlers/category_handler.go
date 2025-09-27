package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"topup-backend/internal/services"
	"topup-backend/internal/utils"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	services  *services.ServiceContainer
	validator *validator.Validate
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(services *services.ServiceContainer) *CategoryHandler {
	return &CategoryHandler{
		services:  services,
		validator: validator.New(),
	}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get all active categories for public view
// @Tags categories
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.Category}
// @Failure 500 {object} utils.APIResponse
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.services.CategoryService.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.APIResponse{data=models.Category}
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// GetCategoriesForAdmin godoc
// @Summary Get categories for admin
// @Description Get paginated categories for admin management
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.APIResponse{data=[]models.Category,pagination=utils.Pagination}
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/categories [get]
func (h *CategoryHandler) GetCategoriesForAdmin(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	categories, total, err := h.services.CategoryService.GetCategoriesForAdmin(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, limit, int(total))
	utils.SuccessResponseWithPagination(c, http.StatusOK, "Categories retrieved successfully", categories, pagination)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category (Admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} utils.APIResponse{data=models.Category}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Router /admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req services.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	category, err := h.services.CategoryService.CreateCategory(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category by ID (Admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body services.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} utils.APIResponse{data=models.Category}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID (Admin only)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}
