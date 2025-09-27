package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"topup-backend/internal/config"
	"topup-backend/internal/models"
)

// Initialize initializes database connection and performs migrations
func Initialize(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Build DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	// Open database connection with safer settings
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Disable foreign key constraint check during migration
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("🗄️  Database connected successfully")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Printf("⚠️  Migration warning: %v", err)
		log.Println("🔄 Database migrations completed with warnings")
	} else {
		log.Println("🔄 Database migrations completed")
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("🔍 Database indexes created")

	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	// Configure GORM to be more lenient with migrations
	migrator := db.Migrator()

	modelsToMigrate := []interface{}{
		&models.User{},
		&models.AdminUser{},
		&models.Category{},
		&models.Product{},
		&models.GameAccount{},
		&models.Voucher{},
		&models.VoucherApplication{},
		&models.Transaction{},
		&models.VoucherUsage{},
		&models.TransactionLog{},
	}

	for _, model := range modelsToMigrate {
		// Check if table exists first
		if !migrator.HasTable(model) {
			log.Printf("📋 Creating new table for %T", model)
			if err := db.AutoMigrate(model); err != nil {
				log.Printf("⚠️  Failed to create table for %T: %v", model, err)
				return err
			}
		} else {
			// Table exists, try to migrate carefully
			log.Printf("🔄 Migrating existing table for %T", model)
			if err := db.AutoMigrate(model); err != nil {
				// If migration fails, log warning but continue
				log.Printf("⚠️  Migration warning for %T: %v - continuing anyway", model, err)
			}
		}
	}

	return nil
}

// createIndexes creates additional database indexes for performance
func createIndexes(db *gorm.DB) error {
	// User indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)").Error; err != nil {
		return err
	}

	// Admin users indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_admin_users_username ON admin_users(username)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_admin_users_role_status ON admin_users(role, status)").Error; err != nil {
		return err
	}

	// Categories indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_status_order ON categories(status, display_order)").Error; err != nil {
		return err
	}

	// Products indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_name ON products USING gin(to_tsvector('english', name))").Error; err != nil {
		log.Println("Warning: Could not create full-text search index for products (PostgreSQL extension might be missing)")
	}

	// Game accounts indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_game_accounts_game_id ON game_accounts(game_id)").Error; err != nil {
		return err
	}

	// Vouchers indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_vouchers_type_status ON vouchers(type, status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_vouchers_application_type ON vouchers(application_type)").Error; err != nil {
		return err
	}

	// Transactions indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_payment_status ON transactions(payment_status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_payment_method ON transactions(payment_method)").Error; err != nil {
		return err
	}

	// Transaction logs indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transaction_logs_created_at ON transaction_logs(created_at DESC)").Error; err != nil {
		return err
	}

	return nil
}

// SeedData seeds initial data for development
func SeedData(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("🌱 Database already seeded, skipping...")
		return nil
	}

	log.Println("🌱 Seeding database with initial data...")

	// Seed categories
	categories := []models.Category{
		{
			Name:         "Mobile Legends",
			Slug:         "mobile-legends",
			Description:  stringPtr("Mobile Legends: Bang Bang Diamond Topup"),
			DisplayOrder: 1,
			Status:       "active",
		},
		{
			Name:         "Free Fire",
			Slug:         "free-fire",
			Description:  stringPtr("Free Fire Diamond Topup"),
			DisplayOrder: 2,
			Status:       "active",
		},
		{
			Name:         "Genshin Impact",
			Slug:         "genshin-impact",
			Description:  stringPtr("Genshin Impact Genesis Crystal Topup"),
			DisplayOrder: 3,
			Status:       "active",
		},
		{
			Name:         "PUBG Mobile",
			Slug:         "pubg-mobile",
			Description:  stringPtr("PUBG Mobile UC Topup"),
			DisplayOrder: 4,
			Status:       "active",
		},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			return err
		}
	}

	// Seed products
	products := []models.Product{
		{
			CategoryID:         1, // Mobile Legends
			Name:               "86 Diamond Mobile Legends",
			Slug:               "86-diamond-mobile-legends",
			Description:        stringPtr("86 Diamond Mobile Legends + 2 Bonus"),
			Price:              20000,
			OriginalPrice:      floatPtr(22000),
			DiscountPercentage: 9,
			Status:             "active",
			DisplayOrder:       1,
			FormConfig: []models.FormField{
				{
					Field:       "gameAccount",
					Label:       "User ID",
					Type:        "text",
					Required:    true,
					Placeholder: "Enter your User ID",
				},
				{
					Field:       "gameZone",
					Label:       "Zone ID",
					Type:        "text",
					Required:    true,
					Placeholder: "Enter your Zone ID",
				},
			},
		},
		{
			CategoryID:   2, // Free Fire
			Name:         "70 Diamond Free Fire",
			Slug:         "70-diamond-free-fire",
			Description:  stringPtr("70 Diamond Free Fire + 7 Bonus"),
			Price:        10000,
			Status:       "active",
			DisplayOrder: 1,
			FormConfig: []models.FormField{
				{
					Field:       "gameAccount",
					Label:       "Player ID",
					Type:        "text",
					Required:    true,
					Placeholder: "Enter your Player ID",
				},
			},
		},
		{
			CategoryID:   3, // Genshin Impact
			Name:         "60 Genesis Crystal",
			Slug:         "60-genesis-crystal-genshin-impact",
			Description:  stringPtr("60 Genesis Crystal Genshin Impact"),
			Price:        15000,
			Status:       "active",
			DisplayOrder: 1,
			FormConfig: []models.FormField{
				{
					Field:       "gameAccount",
					Label:       "UID",
					Type:        "text",
					Required:    true,
					Placeholder: "Enter your UID",
				},
				{
					Field:    "gameServer",
					Label:    "Server",
					Type:     "select",
					Required: true,
					Options:  []string{"Asia", "America", "Europe", "TW/HK/MO"},
				},
			},
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	// Seed sample vouchers
	vouchers := []models.Voucher{
		{
			Code:                 "NEWUSER10",
			Type:                 "percentage",
			Value:                10,
			Description:          stringPtr("10% discount for new users"),
			ApplicationType:      "all",
			MinTransactionAmount: 15000,
			MaxDiscountAmount:    floatPtr(5000),
			Quota:                100,
			MaxUsesPerUser:       1,
			StartDate:            time.Now().AddDate(0, 0, -1), // Yesterday
			EndDate:              time.Now().AddDate(0, 1, 0),  // 1 month from now
			Status:               "active",
		},
		{
			Code:                 "WEEKEND15",
			Type:                 "percentage",
			Value:                15,
			Description:          stringPtr("15% weekend special discount"),
			ApplicationType:      "category",
			MinTransactionAmount: 20000,
			MaxDiscountAmount:    floatPtr(10000),
			Quota:                50,
			MaxUsesPerUser:       2,
			StartDate:            time.Now().AddDate(0, 0, -1),
			EndDate:              time.Now().AddDate(0, 0, 30),
			Status:               "active",
		},
	}

	for _, voucher := range vouchers {
		if err := db.Create(&voucher).Error; err != nil {
			return err
		}
	}

	// Create default admin user
	adminUser := models.AdminUser{
		Username:     "admin",
		Name:         "Super Administrator",
		Email:        stringPtr("admin@wawstore.com"),
		PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: "password123"
		Role:         "super_admin",
		Status:       "active",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	log.Println("✅ Database seeded successfully")
	return nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
