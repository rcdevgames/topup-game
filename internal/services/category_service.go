package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"topup-backend/internal/models"
	"topup-backend/internal/utils"
)

// CategoryService handles category-related business logic
type CategoryService struct {
	db    *gorm.DB
	cache *CacheService
}

// NewCategoryService creates a new category service
func NewCategoryService(db *gorm.DB, cache *CacheService) *CategoryService {
	return &CategoryService{
		db:    db,
		cache: cache,
	}
}

// CreateCategoryRequest represents create category request
type CreateCategoryRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

// UpdateCategoryRequest represents update category request
type UpdateCategoryRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"display_order"`
	Status       string `json:"status" validate:"required,oneof=active inactive"`
}

// GetAllCategories gets all active categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	// Try to get from cache first
	cacheKey := KeyCategories
	if err := s.cache.Get(cacheKey, &categories); err == nil {
		return categories, nil
	}

	// Get from database
	if err := s.db.Where("status = ?", "active").
		Order("display_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, categories, GetCacheTTL(cacheKey))

	return categories, nil
}

// GetCategoryByID gets category by ID
func (s *CategoryService) GetCategoryByID(categoryID uint) (*models.Category, error) {
	var category models.Category

	// Try to get from cache first
	cacheKey := fmt.Sprintf(KeyCategoryByID, categoryID)
	if err := s.cache.Get(cacheKey, &category); err == nil {
		return &category, nil
	}

	// Get from database
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, category, GetCacheTTL(cacheKey))

	return &category, nil
}

// GetCategoriesWithProducts gets categories with their products
func (s *CategoryService) GetCategoriesWithProducts() ([]models.Category, error) {
	var categories []models.Category

	if err := s.db.Preload("Products", "status = ?", "active").
		Where("status = ?", "active").
		Order("display_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories with products: %w", err)
	}

	return categories, nil
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req *CreateCategoryRequest) (*models.Category, error) {
	// Generate slug
	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	var existingCategory models.Category
	if err := s.db.Where("slug = ?", slug).First(&existingCategory).Error; err == nil {
		return nil, errors.New("category with similar name already exists")
	}

	// Create category
	category := models.Category{
		Name:         req.Name,
		Slug:         slug,
		DisplayOrder: req.DisplayOrder,
		Status:       "active",
	}

	if req.Description != "" {
		category.Description = &req.Description
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Clear cache
	s.clearCategoryCache()

	return &category, nil
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(categoryID uint, req *UpdateCategoryRequest) (*models.Category, error) {
	// Find category
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	// Generate new slug if name changed
	if category.Name != req.Name {
		newSlug := utils.GenerateSlug(req.Name)

		// Check if new slug already exists (excluding current category)
		var existingCategory models.Category
		if err := s.db.Where("slug = ? AND id != ?", newSlug, categoryID).First(&existingCategory).Error; err == nil {
			return nil, errors.New("category with similar name already exists")
		}

		category.Slug = newSlug
	}

	// Update fields
	category.Name = req.Name
	category.DisplayOrder = req.DisplayOrder
	category.Status = req.Status

	if req.Description != "" {
		category.Description = &req.Description
	} else {
		category.Description = nil
	}

	// Save changes
	if err := s.db.Save(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Clear cache
	s.clearCategoryCache()

	return &category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(categoryID uint) error {
	// Check if category exists
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return fmt.Errorf("failed to find category: %w", err)
	}

	// Check if category has products
	var productCount int64
	if err := s.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to check products: %w", err)
	}

	if productCount > 0 {
		return errors.New("cannot delete category with existing products")
	}

	// Delete category
	if err := s.db.Delete(&category).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// Clear cache
	s.clearCategoryCache()

	return nil
}

// GetCategoriesForAdmin gets all categories for admin with pagination
func (s *CategoryService) GetCategoriesForAdmin(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	offset := utils.CalculateOffset(page, limit)

	// Count total
	if err := s.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Get categories with pagination
	if err := s.db.Offset(offset).
		Limit(limit).
		Order("display_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, total, nil
}

// UpdateCategoryIcon updates category icon URL
func (s *CategoryService) UpdateCategoryIcon(categoryID uint, iconURL string) (*models.Category, error) {
	// Find category
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	// Update icon URL
	category.IconURL = &iconURL

	// Save changes
	if err := s.db.Save(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to update category icon: %w", err)
	}

	// Clear cache
	s.clearCategoryCache()

	return &category, nil
}

// clearCategoryCache clears all category-related cache
func (s *CategoryService) clearCategoryCache() {
	s.cache.Clear("categories:*")
	s.cache.Clear("category:*")
}
