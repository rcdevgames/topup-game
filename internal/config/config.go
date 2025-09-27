package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Redis       RedisConfig
	AWS         AWSConfig
	Payment     PaymentConfig
	WhatsApp    WhatsAppConfig
	RateLimit   RateLimitConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port      string
	UploadDir string
	BaseURL   string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// AWSConfig holds AWS S3 configuration
type AWSConfig struct {
	S3Bucket        string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

// PaymentConfig holds payment gateway configuration
type PaymentConfig struct {
	MidtransServerKey string
	MidtransClientKey string
	MidtransEnv       string
}

// WhatsAppConfig holds WhatsApp API configuration
type WhatsAppConfig struct {
	APIKey  string
	PhoneID string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("NODE_ENV", "development"),
		Server: ServerConfig{
			Port:      getEnv("PORT", "8080"),
			UploadDir: getEnv("UPLOAD_DIR", "./cdn"),
			BaseURL:   getEnv("BASE_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "topup_user"),
			Password: getEnv("DB_PASSWORD", "topup_password"),
			Name:     getEnv("DB_NAME", "topup_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			AccessExpire:  parseDuration(getEnv("JWT_ACCESS_EXPIRE", "15m")),
			RefreshExpire: parseDuration(getEnv("JWT_REFRESH_EXPIRE", "168h")),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt(getEnv("REDIS_DB", "0")),
		},
		AWS: AWSConfig{
			S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", "ap-southeast-1"),
		},
		Payment: PaymentConfig{
			MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
			MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
			MidtransEnv:       getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
		},
		WhatsApp: WhatsAppConfig{
			APIKey:  getEnv("WHATSAPP_API_KEY", ""),
			PhoneID: getEnv("WHATSAPP_PHONE_ID", ""),
		},
		RateLimit: RateLimitConfig{
			Enabled:           getBool(getEnv("RATE_LIMIT_ENABLED", "true")),
			RequestsPerMinute: parseInt(getEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", "60")),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute // Default fallback
	}
	return duration
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func getBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}
