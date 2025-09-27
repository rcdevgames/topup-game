package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		body           string
		expectedStatus int
		skipCSRF       bool
	}{
		{
			name:           "GET request should pass without CSRF token",
			method:         "GET",
			path:           "/api/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST to auth/login should pass without CSRF token",
			method:         "POST",
			path:           "/api/auth/login",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST to auth/register should pass without CSRF token",
			method:         "POST",
			path:           "/api/auth/register",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET CSRF token should pass without CSRF token",
			method:         "GET",
			path:           "/api/csrf",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST without CSRF token should fail",
			method:         "POST",
			path:           "/api/protected",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "POST with invalid CSRF token should fail",
			method:         "POST",
			path:           "/api/protected",
			headers:        map[string]string{"X-CSRF-Token": "invalid-token"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "PUT without CSRF token should fail",
			method:         "PUT",
			path:           "/api/protected",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "PATCH without CSRF token should fail",
			method:         "PATCH",
			path:           "/api/protected",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "DELETE without CSRF token should fail",
			method:         "DELETE",
			path:           "/api/protected",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(CSRFMiddleware())
			router.Use(CSRFCheck())

			// Add test routes
			router.GET("/api/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})
			router.POST("/api/auth/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "login success"})
			})
			router.POST("/api/auth/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "register success"})
			})
			router.GET("/api/csrf", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "csrf token"})
			})
			router.POST("/api/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "protected success"})
			})
			router.PUT("/api/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "protected success"})
			})
			router.PATCH("/api/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "protected success"})
			})
			router.DELETE("/api/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "protected success"})
			})

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			// Add headers
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCSRFTokenGeneration(t *testing.T) {
	// Test token generation
	token1 := generateCSRFToken()
	token2 := generateCSRFToken()

	// Tokens should not be empty
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)

	// Tokens should be different (random)
	assert.NotEqual(t, token1, token2)

	// Tokens should contain a dot (signature separator)
	assert.Contains(t, token1, ".")
	assert.Contains(t, token2, ".")
}

func TestCSRFTokenValidation(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "Valid token should pass",
			token:    generateCSRFToken(),
			expected: true,
		},
		{
			name:     "Empty token should fail",
			token:    "",
			expected: false,
		},
		{
			name:     "Invalid format token should fail",
			token:    "invalid-token-format",
			expected: false,
		},
		{
			name:     "Token without dot should fail",
			token:    "invalidtokenwithoutdot",
			expected: false,
		},
		{
			name:     "Token with invalid base64 should fail",
			token:    "invalid.base64!@#",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateCSRFToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCSRFSecret(t *testing.T) {
	// Should return default secret when env var not set
	secret := getCSRFSecret()
	assert.NotEmpty(t, secret)
}

func TestCSRFTokenFromDifferentSources(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CSRFMiddleware())

	// Generate a valid token
	validToken := generateCSRFToken()

	// Add test route that requires CSRF
	router.POST("/test", CSRFCheck(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name           string
		tokenSource    string // "header", "form", "query"
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token in header",
			tokenSource:    "header",
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid token in form",
			tokenSource:    "form",
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid token in query",
			tokenSource:    "query",
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			switch tt.tokenSource {
			case "header":
				req = httptest.NewRequest("POST", "/test", nil)
				req.Header.Set("X-CSRF-Token", tt.token)
			case "form":
				body := "csrf_token=" + tt.token
				req = httptest.NewRequest("POST", "/test", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			case "query":
				req = httptest.NewRequest("POST", "/test?csrf_token="+tt.token, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
