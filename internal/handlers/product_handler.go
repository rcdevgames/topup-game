package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"topup-backend/internal/services"
	"topup-backend/internal/utils"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	services  *services.ServiceContainer
	validator *validator.Validate
}

// NewProductHandler creates a new product handler
func NewProductHandler(services *services.ServiceContainer) *ProductHandler {
	return &ProductHandler{
		services:  services,
		validator: validator.New(),
	}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all active products for public view
// @Tags products
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.Product}
// @Failure 500 {object} utils.APIResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.services.ProductService.GetAllProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a specific product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.APIResponse{data=models.Product}
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products by name or description
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param category_id query int false "Category ID filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.APIResponse{data=[]models.Product,pagination=utils.Pagination}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get all products in a specific category
// @Tags products
// @Produce json
// @Param category_id path int true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.APIResponse{data=[]models.Product,pagination=utils.Pagination}
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// GetProductsForAdmin godoc
// @Summary Get products for admin
// @Description Get paginated list of products for admin management
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(active, inactive, out_of_stock)
// @Param category_id query int false "Filter by category ID"
// @Success 200 {object} utils.APIResponse{data=[]models.Product,pagination=utils.Pagination}
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/products [get]
func (h *ProductHandler) GetProductsForAdmin(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	products, total, err := h.services.ProductService.GetProductsForAdmin(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, limit, int(total))
	utils.SuccessResponseWithPagination(c, http.StatusOK, "Products retrieved successfully", products, pagination)
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product (Admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.CreateProductRequest true "Product creation request"
// @Success 201 {object} utils.APIResponse{data=models.Product}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req services.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	product, err := h.services.ProductService.CreateProduct(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

// GetProductByIDForAdmin godoc
// @Summary Get product by ID for admin
// @Description Get detailed product information by ID for admin
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.APIResponse{data=models.Product}
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /admin/products/{id} [get]
func (h *ProductHandler) GetProductByIDForAdmin(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented yet")
}
