package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"topup-backend/internal/config"
)

// APIResponse represents standard API response format
type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithPagination sends a success response with pagination
func SuccessResponseWithPagination(c *gin.Context, statusCode int, message string, data interface{}, pagination *Pagination) {
	c.JSON(statusCode, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}

// ValidationErrorResponse sends validation error response
func ValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates JWT tokens
func GenerateJWT(userID uint, userType, role string, cfg *config.Config) (accessToken, refreshToken string, err error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"type": userType,
		"role": role,
		"exp":  time.Now().Add(cfg.JWT.AccessExpire).Unix(),
		"iat":  time.Now().Unix(),
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"exp":  time.Now().Add(cfg.JWT.RefreshExpire).Unix(),
		"iat":  time.Now().Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseRefreshToken parses and validates refresh token
func ParseRefreshToken(tokenString string, cfg *config.Config) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	// Check if it's a refresh token
	if claims["type"] != "refresh" {
		return 0, jwt.ErrInvalidType
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return uint(userID), nil
}

// GenerateSlug generates URL-friendly slug from title
func GenerateSlug(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// GetPaginationParams extracts pagination parameters from query
func GetPaginationParams(c *gin.Context) (page, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	return page, limit
}

// CalculateOffset calculates database offset from page and limit
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculatePagination calculates pagination info
func CalculatePagination(page, limit, total int) *Pagination {
	totalPages := (total + limit - 1) / limit
	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// IsValidEmail checks if email format is valid
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// IsValidPhone checks if phone number format is valid (Indonesian format)
func IsValidPhone(phone string) bool {
	// Remove all non-numeric characters
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// Check length (Indonesian phone numbers are typically 10-15 digits)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return false
	}

	// Check if all characters are digits
	for _, r := range cleaned {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Indonesian phone numbers typically start with 08, 62, or 8
	return strings.HasPrefix(cleaned, "08") || strings.HasPrefix(cleaned, "62") || strings.HasPrefix(cleaned, "8")
}

// NormalizePhone normalizes phone number to Indonesian format
func NormalizePhone(phone string) string {
	// Remove all non-numeric characters
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// Convert to standard format starting with 62
	if strings.HasPrefix(cleaned, "08") {
		cleaned = "62" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "8") {
		cleaned = "62" + cleaned
	}

	return cleaned
}

// StringPtr returns a pointer to string
func StringPtr(s string) *string {
	return &s
}

// UintPtr returns a pointer to uint
func UintPtr(u uint) *uint {
	return &u
}

// Float64Ptr returns a pointer to float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// GetClientIP gets client IP address from gin context
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to ClientIP from gin
	return c.ClientIP()
}
