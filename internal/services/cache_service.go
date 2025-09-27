package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"topup-backend/internal/config"
)

// CacheService handles caching operations
type CacheService struct {
	client *redis.Client
	cfg    *config.Config
}

// NewCacheService creates a new cache service
func NewCacheService(cfg *config.Config) *CacheService {
	// If Redis URL is not configured, return a no-op cache service
	if cfg.Redis.URL == "" {
		log.Println("⚠️  Redis not configured, caching disabled")
		return &CacheService{client: nil, cfg: cfg}
	}

	// Parse Redis URL and create client
	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Printf("⚠️  Failed to parse Redis URL: %v, caching disabled", err)
		return &CacheService{client: nil, cfg: cfg}
	}

	if cfg.Redis.Password != "" {
		opts.Password = cfg.Redis.Password
	}
	opts.DB = cfg.Redis.DB

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️  Failed to connect to Redis: %v, caching disabled", err)
		return &CacheService{client: nil, cfg: cfg}
	}

	log.Println("🚀 Redis connected successfully")
	return &CacheService{client: client, cfg: cfg}
}

// Set stores a value in cache with TTL
func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	if s.client == nil {
		return nil // No-op if Redis is not available
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (s *CacheService) Get(key string, dest interface{}) error {
	if s.client == nil {
		return redis.Nil // Return cache miss if Redis is not available
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key from cache
func (s *CacheService) Delete(key string) error {
	if s.client == nil {
		return nil // No-op if Redis is not available
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.Del(ctx, key).Err()
}

// Clear removes all keys matching a pattern
func (s *CacheService) Clear(pattern string) error {
	if s.client == nil {
		return nil // No-op if Redis is not available
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all keys matching pattern
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all matching keys
	return s.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in cache
func (s *CacheService) Exists(key string) bool {
	if s.client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := s.client.Exists(ctx, key).Result()
	return err == nil && count > 0
}

// Cache key constants and TTL values
var (
	// Cache TTL durations
	CacheShort  = 5 * time.Minute  // For frequently changing data
	CacheMedium = 30 * time.Minute // For moderately changing data
	CacheLong   = 2 * time.Hour    // For rarely changing data
	CacheDaily  = 24 * time.Hour   // For daily data

	// Cache key patterns
	KeyProducts           = "products:active"
	KeyProductByID        = "product:%d"
	KeyCategories         = "categories:active"
	KeyCategoryByID       = "category:%d"
	KeyVoucherValidation  = "voucher:%s"
	KeyUserProfile        = "user:profile:%d"
	KeyDashboardAnalytics = "analytics:dashboard"
	KeyTopProducts        = "analytics:top_products"
	KeyDailySales         = "analytics:daily_sales"
)

// GetCacheTTL returns appropriate TTL for cache key
func GetCacheTTL(key string) time.Duration {
	switch {
	case key == KeyProducts || key == KeyCategories:
		return CacheLong
	case key == KeyDashboardAnalytics:
		return CacheShort
	case key == KeyTopProducts || key == KeyDailySales:
		return CacheMedium
	default:
		return CacheMedium
	}
}
