package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware returns a CSRF protection middleware
func CSRFMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Generate CSRF token for the session and store in context
		token := generateCSRFToken()
		c.Set("csrf_token", token)
		c.Next()
	})
}

// CSRFCheck middleware to validate CSRF token for specific methods
func CSRFCheck() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip CSRF check for safe methods and health check
		method := c.Request.Method
		path := c.Request.URL.Path

		// Skip for safe methods (GET, HEAD, OPTIONS)
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		// Skip for health check and public endpoints
		skipPaths := []string{
			"/health",
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/csrf",
			"/docs",
			"/swagger",
		}

		for _, skipPath := range skipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// Check if CSRF token exists and validate it
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("csrf_token")
		}
		if token == "" {
			token = c.Query("csrf_token")
		}

		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "CSRF token is required for this operation",
				"error":   "csrf_token_missing",
			})
			c.Abort()
			return
		}

		// Validate CSRF token
		if !validateCSRFToken(token) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "CSRF token is invalid or expired",
				"error":   "csrf_token_invalid",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

func getCSRFSecret() string {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		secret = "default-csrf-secret-change-in-production"
	}
	return secret
}

// GetCSRFToken extracts CSRF token from context
func GetCSRFToken(c *gin.Context) string {
	if token := c.GetHeader("X-CSRF-Token"); token != "" {
		return token
	}
	if token := c.PostForm("csrf_token"); token != "" {
		return token
	}
	if token := c.Query("csrf_token"); token != "" {
		return token
	}
	return ""
}

// generateCSRFToken generates a new CSRF token
func generateCSRFToken() string {
	// Create a random byte slice
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp-based token if random generation fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Add timestamp to make token time-sensitive
	timestamp := time.Now().Unix()
	tokenData := fmt.Sprintf("%s:%d", hex.EncodeToString(randomBytes), timestamp)

	// Sign the token with HMAC
	secret := getCSRFSecret()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tokenData))
	signature := h.Sum(nil)

	// Combine token data and signature
	token := fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString([]byte(tokenData)),
		base64.URLEncoding.EncodeToString(signature))

	return token
}

// validateCSRFToken validates a CSRF token
func validateCSRFToken(token string) bool {
	if token == "" {
		return false
	}

	// Split token into data and signature parts
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	// Decode token data
	tokenDataBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	tokenData := string(tokenDataBytes)

	// Decode signature
	expectedSignature, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Verify signature
	secret := getCSRFSecret()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tokenData))
	actualSignature := h.Sum(nil)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return false
	}

	// Check token age (tokens expire after 24 hours)
	parts = strings.Split(tokenData, ":")
	if len(parts) != 2 {
		return false
	}

	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}

	// Token is valid for 24 hours
	tokenAge := time.Now().Unix() - timestamp
	return tokenAge <= 24*60*60 // 24 hours in seconds
}

// CSRFTokenResponse represents the CSRF token response
type CSRFTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
