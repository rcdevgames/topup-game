package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GameAccountData represents denormalized game account data in transactions
type GameAccountData struct {
	GameAccount string `json:"game_account" validate:"required"`
	GameZone    string `json:"game_zone,omitempty"`
	GameServer  string `json:"game_server,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
}

// Transaction represents a topup transaction
type Transaction struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	TransactionCode string `gorm:"size:20;unique;not null;index:idx_transaction_code" json:"transaction_code"`
	UserID          uint   `gorm:"not null;index:idx_user_status" json:"user_id" validate:"required"`
	ProductID       uint   `gorm:"not null" json:"product_id" validate:"required"`

	// Game Account Info (denormalized for transaction history)
	GameAccountData GameAccountData `gorm:"type:jsonb;serializer:json;not null" json:"game_account_data" validate:"required"`

	// Pricing
	ProductPrice    float64 `gorm:"type:decimal(10,2);not null" json:"product_price" validate:"required,min=0"`
	PaymentFee      float64 `gorm:"type:decimal(10,2);default:0" json:"payment_fee" validate:"min=0"`
	VoucherDiscount float64 `gorm:"type:decimal(10,2);default:0" json:"voucher_discount" validate:"min=0"`
	TotalAmount     float64 `gorm:"type:decimal(10,2);not null" json:"total_amount" validate:"required,min=0"`

	// Payment
	PaymentMethod    string  `gorm:"size:50;not null" json:"payment_method" validate:"required,oneof=gopay ovo dana bca mandiri bni"`
	PaymentStatus    string  `gorm:"size:20;default:pending" json:"payment_status" validate:"oneof=pending paid failed expired refunded"`
	PaymentReference *string `gorm:"size:100" json:"payment_reference"`
	PaymentURL       *string `gorm:"type:text" json:"payment_url"`

	// Contact
	WhatsApp string `gorm:"size:15;not null" json:"whatsapp" validate:"required,phone"`

	// Status & Processing
	Status      string     `gorm:"size:20;default:pending;index:idx_user_status,idx_status_created" json:"status" validate:"oneof=pending processing completed failed cancelled"`
	ProcessedAt *time.Time `json:"processed_at"`
	CompletedAt *time.Time `json:"completed_at"`
	ExpiredAt   *time.Time `json:"expired_at"`

	// Metadata
	UserAgent *string `gorm:"type:text" json:"user_agent"`
	IPAddress *string `gorm:"size:45" json:"ip_address"`

	CreatedAt time.Time `gorm:"index:idx_status_created" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User          User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product       Product          `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Logs          []TransactionLog `gorm:"foreignKey:TransactionID" json:"logs,omitempty"`
	VoucherUsages []VoucherUsage   `gorm:"foreignKey:TransactionID" json:"voucher_usages,omitempty"`
}

// TransactionLog represents transaction status change history
type TransactionLog struct {
	ID             uint                   `gorm:"primaryKey" json:"id"`
	TransactionID  uint                   `gorm:"not null" json:"transaction_id" validate:"required"`
	StatusFrom     *string                `gorm:"size:20" json:"status_from"`
	StatusTo       string                 `gorm:"size:20;not null" json:"status_to" validate:"required"`
	Message        *string                `gorm:"type:text" json:"message"`
	Metadata       map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"metadata,omitempty"`
	CreatedByAdmin *uint                  `json:"created_by_admin"`
	CreatedAt      time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Admin       *AdminUser  `gorm:"foreignKey:CreatedByAdmin" json:"admin,omitempty"`
}

// BeforeCreate hook for Transaction
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// Generate unique transaction code
	if t.TransactionCode == "" {
		t.TransactionCode = generateTransactionCode()
	}

	// Set expiry time (24 hours from creation)
	if t.ExpiredAt == nil {
		expiryTime := time.Now().Add(24 * time.Hour)
		t.ExpiredAt = &expiryTime
	}

	return nil
}

// AfterCreate hook for Transaction
func (t *Transaction) AfterCreate(tx *gorm.DB) error {
	// Create initial transaction log
	log := TransactionLog{
		TransactionID: t.ID,
		StatusFrom:    nil,
		StatusTo:      t.Status,
		Message:       stringPtr("Transaction created"),
		CreatedAt:     time.Now(),
	}

	return tx.Create(&log).Error
}

// BeforeUpdate hook for Transaction
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	// Get the original transaction to compare status
	var original Transaction
	if err := tx.First(&original, t.ID).Error; err != nil {
		return err
	}

	// If status changed, create a log entry
	if original.Status != t.Status {
		log := TransactionLog{
			TransactionID: t.ID,
			StatusFrom:    &original.Status,
			StatusTo:      t.Status,
			Message:       stringPtr("Status updated"),
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// Update processed/completed timestamps
		now := time.Now()
		switch t.Status {
		case "processing":
			if t.ProcessedAt == nil {
				t.ProcessedAt = &now
			}
		case "completed":
			if t.CompletedAt == nil {
				t.CompletedAt = &now
			}
		}
	}

	return nil
}

// generateTransactionCode generates a unique transaction code
func generateTransactionCode() string {
	// Generate UUID and take first 8 characters + prefix
	id := uuid.New().String()
	return "TXN" + id[:8]
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
