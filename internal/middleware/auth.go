package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"topup-backend/internal/config"
	"topup-backend/internal/utils"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID uint   `json:"sub"`
	Type   string `json:"type"` // "user" or "admin"
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Check Bearer prefix
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		// Check if token is not a refresh token
		if claims.Type == "refresh" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Cannot use refresh token for authentication")
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.Type)
		c.Set("user_role", claims.Role)

		c.Next()
	})
}

// AdminAuthMiddleware ensures the user is an admin
func AdminAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Run JWT auth first
		AuthMiddleware(cfg)(c)

		if c.IsAborted() {
			return
		}

		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	})
}

// RoleMiddleware checks if admin has required role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "User role not found")
			c.Abort()
			return
		}

		role := userRole.(string)

		// Super admin has access to everything
		if role == "super_admin" {
			c.Next()
			return
		}

		// Check if user has required role
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
		c.Abort()
	})
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow specific origins or all for development
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://wawstore.com",
			"https://admin.wawstore.com",
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// More permissive CSP for Swagger UI to work properly
		// In production, you might want to be more restrictive
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/swagger/") {
			// Allow inline styles and scripts for Swagger UI
			c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")
		} else {
			// More restrictive CSP for other endpoints
			c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self'; font-src 'self'")
		}

		c.Next()
	})
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	if !cfg.RateLimit.Enabled {
		return gin.HandlerFunc(func(c *gin.Context) {
			c.Next()
		})
	}

	// Simple in-memory rate limiter (in production, use Redis)
	requests := make(map[string][]time.Time)

	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old requests
		if times, exists := requests[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			requests[clientIP] = validTimes
		}

		// Check rate limit
		if len(requests[clientIP]) >= cfg.RateLimit.RequestsPerMinute {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}

		// Add current request
		requests[clientIP] = append(requests[clientIP], now)

		c.Next()
	})
}
