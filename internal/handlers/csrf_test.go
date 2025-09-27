package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCSRFHandler_GetCSRFToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCSRFHandler()
	router := gin.New()

	// Mock CSRF middleware that sets token in context
	router.Use(func(c *gin.Context) {
		c.Set("csrf_token", "test-csrf-token-12345")
		c.Next()
	})

	router.GET("/csrf", handler.GetCSRFToken)

	req := httptest.NewRequest("GET", "/csrf", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "CSRF token generated successfully")
	assert.Contains(t, w.Body.String(), "test-csrf-token-12345")
	assert.Contains(t, w.Body.String(), "X-CSRF-Token")
	assert.Contains(t, w.Body.String(), "expires_at")
}

func TestCSRFHandler_GetCSRFToken_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCSRFHandler()
	router := gin.New()

	// No middleware to set token in context
	router.GET("/csrf", handler.GetCSRFToken)

	req := httptest.NewRequest("GET", "/csrf", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to generate CSRF token")
	assert.Contains(t, w.Body.String(), "csrf_generation_failed")
}

func TestCSRFHandler_ValidateCSRF(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCSRFHandler()
	router := gin.New()

	// This endpoint assumes CSRF validation is handled by middleware
	router.POST("/csrf/validate", handler.ValidateCSRF)

	req := httptest.NewRequest("POST", "/csrf/validate", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "CSRF token is valid")
	assert.Contains(t, w.Body.String(), "token_valid")
}
