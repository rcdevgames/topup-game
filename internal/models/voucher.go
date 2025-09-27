package models

import (
	"time"

	"gorm.io/gorm"
)

// Voucher represents a discount voucher
type Voucher struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Code                 string    `gorm:"size:50;unique;not null;index:idx_code_status" json:"code" validate:"required,uppercase,min=3,max=50"`
	Type                 string    `gorm:"size:20;not null" json:"type" validate:"required,oneof=percentage fixed"`
	Value                float64   `gorm:"type:decimal(10,2);not null" json:"value" validate:"required,min=0"`
	Description          *string   `gorm:"type:text" json:"description"`
	ApplicationType      string    `gorm:"size:20;default:all" json:"application_type" validate:"oneof=all category product"`
	MinTransactionAmount float64   `gorm:"type:decimal(10,2);default:0" json:"min_transaction_amount" validate:"min=0"`
	MaxDiscountAmount    *float64  `gorm:"type:decimal(10,2)" json:"max_discount_amount" validate:"omitempty,min=0"`
	Quota                int       `gorm:"not null" json:"quota" validate:"required,min=1"`
	UsedCount            int       `gorm:"default:0" json:"used_count"`
	MaxUsesPerUser       int       `gorm:"default:1" json:"max_uses_per_user" validate:"min=1"`
	StartDate            time.Time `gorm:"type:date;not null;index:idx_dates" json:"start_date" validate:"required"`
	EndDate              time.Time `gorm:"type:date;not null;index:idx_dates" json:"end_date" validate:"required,gtfield=StartDate"`
	Status               string    `gorm:"size:20;default:active;index:idx_code_status" json:"status" validate:"oneof=active inactive expired"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relations
	Applications []VoucherApplication `gorm:"foreignKey:VoucherID" json:"applications,omitempty"`
	Usages       []VoucherUsage       `gorm:"foreignKey:VoucherID" json:"usages,omitempty"`
}

// VoucherApplication represents voucher applicability to categories or products
type VoucherApplication struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	VoucherID      uint      `gorm:"not null" json:"voucher_id" validate:"required"`
	ApplicableID   uint      `gorm:"not null" json:"applicable_id" validate:"required"`
	ApplicableType string    `gorm:"size:20;not null" json:"applicable_type" validate:"required,oneof=category product"`
	CreatedAt      time.Time `json:"created_at"`

	// Relations
	Voucher Voucher `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}

// VoucherUsage represents voucher usage history
type VoucherUsage struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	VoucherID      uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"voucher_id" validate:"required"`
	UserID         uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"user_id" validate:"required"`
	TransactionID  uint      `gorm:"not null;uniqueIndex:idx_unique_usage" json:"transaction_id" validate:"required"`
	DiscountAmount float64   `gorm:"type:decimal(10,2);not null" json:"discount_amount" validate:"required,min=0"`
	UsedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"used_at"`

	// Relations
	Voucher     Voucher     `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}

// BeforeCreate hook for Voucher to validate dates
func (v *Voucher) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if v.StartDate.After(v.EndDate) {
		return gorm.ErrInvalidData
	}
	if v.EndDate.Before(now.Truncate(24 * time.Hour)) {
		v.Status = "expired"
	}
	return nil
}

// BeforeUpdate hook for Voucher to validate dates and status
func (v *Voucher) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	if v.StartDate.After(v.EndDate) {
		return gorm.ErrInvalidData
	}
	if v.EndDate.Before(now.Truncate(24 * time.Hour)) {
		v.Status = "expired"
	}
	if v.UsedCount >= v.Quota {
		v.Status = "inactive"
	}
	return nil
}
