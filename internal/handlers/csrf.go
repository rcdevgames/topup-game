package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFHandler handles CSRF related operations
type CSRFHandler struct{}

// NewCSRFHandler creates a new CSRF handler
func NewCSRFHandler() *CSRFHandler {
	return &CSRFHandler{}
}

// GetCSRFToken godoc
// @Summary Get CSRF token
// @Description Get CSRF token for form submission protection
// @Tags CSRF
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "success"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /api/csrf [get]
func (h *CSRFHandler) GetCSRFToken(c *gin.Context) {
	// Get token from context (set by middleware)
	token, exists := c.Get("csrf_token")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate CSRF token",
			"error":   "csrf_generation_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CSRF token generated successfully",
		"data": gin.H{
			"token":      token.(string),
			"expires_at": time.Now().Add(24 * time.Hour), // Token expires in 24 hours
			"usage": gin.H{
				"header":      "X-CSRF-Token",
				"form_field":  "csrf_token",
				"query_param": "csrf_token",
			},
		},
	})
}

// ValidateCSRF godoc
// @Summary Validate CSRF token
// @Description Validate if provided CSRF token is valid
// @Tags CSRF
// @Accept json
// @Produce json
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 403 {object} map[string]interface{} "error"
// @Router /api/v1/csrf/validate [post]
func (h *CSRFHandler) ValidateCSRF(c *gin.Context) {
	// This endpoint is mainly for testing CSRF validation
	// The actual validation is handled by the middleware

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CSRF token is valid",
		"data": gin.H{
			"token_valid": true,
			"timestamp":   time.Now(),
		},
	})
}
