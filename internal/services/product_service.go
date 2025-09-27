package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"topup-backend/internal/models"
	"topup-backend/internal/utils"
)

// ProductService handles product-related business logic
type ProductService struct {
	db    *gorm.DB
	cache *CacheService
}

// NewProductService creates a new product service
func NewProductService(db *gorm.DB, cache *CacheService) *ProductService {
	return &ProductService{
		db:    db,
		cache: cache,
	}
}

// CreateProductRequest represents create product request
type CreateProductRequest struct {
	CategoryID         uint               `json:"category_id" validate:"required"`
	Name               string             `json:"name" validate:"required,min=2,max=200"`
	Description        string             `json:"description,omitempty"`
	Price              float64            `json:"price" validate:"required,min=0"`
	OriginalPrice      *float64           `json:"original_price,omitempty" validate:"omitempty,min=0"`
	DiscountPercentage int                `json:"discount_percentage" validate:"min=0,max=100"`
	FormConfig         []models.FormField `json:"form_config" validate:"required,dive"`
	DisplayOrder       int                `json:"display_order"`
}

// UpdateProductRequest represents update product request
type UpdateProductRequest struct {
	CategoryID         uint               `json:"category_id" validate:"required"`
	Name               string             `json:"name" validate:"required,min=2,max=200"`
	Description        string             `json:"description,omitempty"`
	Price              float64            `json:"price" validate:"required,min=0"`
	OriginalPrice      *float64           `json:"original_price,omitempty" validate:"omitempty,min=0"`
	DiscountPercentage int                `json:"discount_percentage" validate:"min=0,max=100"`
	FormConfig         []models.FormField `json:"form_config" validate:"required,dive"`
	Status             string             `json:"status" validate:"required,oneof=active inactive out_of_stock"`
	DisplayOrder       int                `json:"display_order"`
}

// ProductSearchRequest represents product search parameters
type ProductSearchRequest struct {
	Query      string `json:"query,omitempty"`
	CategoryID *uint  `json:"category_id,omitempty"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// GetAllProducts gets all active products
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product

	// Try to get from cache first
	cacheKey := KeyProducts
	if err := s.cache.Get(cacheKey, &products); err == nil {
		return products, nil
	}

	// Get from database
	if err := s.db.Preload("Category").
		Where("status = ?", "active").
		Order("display_order ASC, name ASC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, products, GetCacheTTL(cacheKey))

	return products, nil
}

// GetProductByID gets product by ID
func (s *ProductService) GetProductByID(productID uint) (*models.Product, error) {
	var product models.Product

	// Try to get from cache first
	cacheKey := fmt.Sprintf(KeyProductByID, productID)
	if err := s.cache.Get(cacheKey, &product); err == nil {
		return &product, nil
	}

	// Get from database
	if err := s.db.Preload("Category").First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, product, GetCacheTTL(cacheKey))

	return &product, nil
}

// GetProductsByCategory gets products by category ID
func (s *ProductService) GetProductsByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product

	if err := s.db.Preload("Category").
		Where("category_id = ? AND status = ?", categoryID, "active").
		Order("display_order ASC, name ASC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	return products, nil
}

// SearchProducts searches products with filters
func (s *ProductService) SearchProducts(req *ProductSearchRequest) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).Preload("Category")

	// Apply filters
	if req.Query != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// Only active products for public search
	query = query.Where("status = ?", "active")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply pagination
	offset := utils.CalculateOffset(req.Page, req.Limit)
	if err := query.Offset(offset).
		Limit(req.Limit).
		Order("display_order ASC, name ASC").
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req *CreateProductRequest) (*models.Product, error) {
	// Verify category exists
	var category models.Category
	if err := s.db.First(&category, req.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to verify category: %w", err)
	}

	// Generate slug
	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	var existingProduct models.Product
	if err := s.db.Where("slug = ?", slug).First(&existingProduct).Error; err == nil {
		return nil, errors.New("product with similar name already exists")
	}

	// Validate form config
	if err := s.validateFormConfig(req.FormConfig); err != nil {
		return nil, fmt.Errorf("invalid form config: %w", err)
	}

	// Create product
	product := models.Product{
		CategoryID:         req.CategoryID,
		Name:               req.Name,
		Slug:               slug,
		Price:              req.Price,
		OriginalPrice:      req.OriginalPrice,
		DiscountPercentage: req.DiscountPercentage,
		FormConfig:         req.FormConfig,
		Status:             "active",
		DisplayOrder:       req.DisplayOrder,
	}

	if req.Description != "" {
		product.Description = &req.Description
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Load category relationship
	if err := s.db.Preload("Category").First(&product, product.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load product with category: %w", err)
	}

	// Clear cache
	s.clearProductCache()

	return &product, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(productID uint, req *UpdateProductRequest) (*models.Product, error) {
	// Find product
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	// Verify category exists
	var category models.Category
	if err := s.db.First(&category, req.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to verify category: %w", err)
	}

	// Generate new slug if name changed
	if product.Name != req.Name {
		newSlug := utils.GenerateSlug(req.Name)

		// Check if new slug already exists (excluding current product)
		var existingProduct models.Product
		if err := s.db.Where("slug = ? AND id != ?", newSlug, productID).First(&existingProduct).Error; err == nil {
			return nil, errors.New("product with similar name already exists")
		}

		product.Slug = newSlug
	}

	// Validate form config
	if err := s.validateFormConfig(req.FormConfig); err != nil {
		return nil, fmt.Errorf("invalid form config: %w", err)
	}

	// Update fields
	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Price = req.Price
	product.OriginalPrice = req.OriginalPrice
	product.DiscountPercentage = req.DiscountPercentage
	product.FormConfig = req.FormConfig
	product.Status = req.Status
	product.DisplayOrder = req.DisplayOrder

	if req.Description != "" {
		product.Description = &req.Description
	} else {
		product.Description = nil
	}

	// Save changes
	if err := s.db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Load category relationship
	if err := s.db.Preload("Category").First(&product, product.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load product with category: %w", err)
	}

	// Clear cache
	s.clearProductCache()

	return &product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(productID uint) error {
	// Check if product exists
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to find product: %w", err)
	}

	// Check if product has transactions
	var transactionCount int64
	if err := s.db.Model(&models.Transaction{}).Where("product_id = ?", productID).Count(&transactionCount).Error; err != nil {
		return fmt.Errorf("failed to check transactions: %w", err)
	}

	if transactionCount > 0 {
		return errors.New("cannot delete product with existing transactions")
	}

	// Delete product
	if err := s.db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Clear cache
	s.clearProductCache()

	return nil
}

// GetProductsForAdmin gets all products for admin with pagination
func (s *ProductService) GetProductsForAdmin(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := utils.CalculateOffset(page, limit)

	// Count total
	if err := s.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get products with pagination
	if err := s.db.Preload("Category").
		Offset(offset).
		Limit(limit).
		Order("display_order ASC, name ASC").
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}

	return products, total, nil
}

// UpdateProductImage updates product image URL
func (s *ProductService) UpdateProductImage(productID uint, imageURL string) (*models.Product, error) {
	// Find product
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	// Update image URL
	product.ImageURL = &imageURL

	// Save changes
	if err := s.db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product image: %w", err)
	}

	// Load category relationship
	if err := s.db.Preload("Category").First(&product, product.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load product with category: %w", err)
	}

	// Clear cache
	s.clearProductCache()

	return &product, nil
}

// validateFormConfig validates product form configuration
func (s *ProductService) validateFormConfig(formConfig []models.FormField) error {
	if len(formConfig) == 0 {
		return errors.New("form config cannot be empty")
	}

	requiredFields := map[string]bool{}
	for _, field := range formConfig {
		// Check required fields
		if field.Field == "" || field.Label == "" || field.Type == "" {
			return errors.New("field, label, and type are required for all form fields")
		}

		// Check for duplicate fields
		if requiredFields[field.Field] {
			return fmt.Errorf("duplicate field: %s", field.Field)
		}
		requiredFields[field.Field] = true

		// Validate field types
		validTypes := map[string]bool{
			"text":   true,
			"select": true,
			"number": true,
		}
		if !validTypes[field.Type] {
			return fmt.Errorf("invalid field type: %s", field.Type)
		}

		// For select type, options are required
		if field.Type == "select" && len(field.Options) == 0 {
			return fmt.Errorf("options are required for select field: %s", field.Field)
		}
	}

	// Ensure gameAccount field exists
	if !requiredFields["gameAccount"] {
		return errors.New("gameAccount field is required in form config")
	}

	return nil
}

// clearProductCache clears all product-related cache
func (s *ProductService) clearProductCache() {
	s.cache.Clear("products:*")
	s.cache.Clear("product:*")
}
