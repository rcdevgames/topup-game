package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a customer user
type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Phone           string     `gorm:"size:15;unique;not null" json:"phone" validate:"required,phone"`
	PasswordHash    string     `gorm:"size:255;not null" json:"-"`
	Email           *string    `gorm:"size:100" json:"email" validate:"omitempty,email"`
	Status          string     `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive suspended"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations
	GameAccounts []GameAccount `gorm:"foreignKey:UserID" json:"game_accounts,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

// AdminUser represents an admin user
type AdminUser struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:50;unique;not null" json:"username" validate:"required,min=3,max=50"`
	Name         string     `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Email        *string    `gorm:"size:100" json:"email" validate:"omitempty,email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         string     `gorm:"size:20;default:operator" json:"role" validate:"oneof=super_admin admin operator moderator"`
	Status       string     `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedBy    *uint      `gorm:"index" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Creator *AdminUser `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// Category represents a product category
type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Slug         string    `gorm:"size:100;unique;not null" json:"slug"`
	Description  *string   `gorm:"type:text" json:"description"`
	IconURL      *string   `gorm:"size:255" json:"icon_url"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	Status       string    `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// FormField represents dynamic form configuration for products
type FormField struct {
	Field       string   `json:"field" validate:"required"`
	Label       string   `json:"label" validate:"required"`
	Type        string   `json:"type" validate:"required,oneof=text select number"`
	Required    bool     `json:"required"`
	Placeholder string   `json:"placeholder,omitempty"`
	Options     []string `json:"options,omitempty"`
}

// Product represents a game topup product
type Product struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	CategoryID         uint        `gorm:"not null;index:idx_category_status" json:"category_id" validate:"required"`
	Name               string      `gorm:"size:200;not null" json:"name" validate:"required,min=2,max=200"`
	Slug               string      `gorm:"size:200;unique;not null" json:"slug"`
	Description        *string     `gorm:"type:text" json:"description"`
	Price              float64     `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
	OriginalPrice      *float64    `gorm:"type:decimal(10,2)" json:"original_price" validate:"omitempty,min=0"`
	DiscountPercentage int         `gorm:"default:0" json:"discount_percentage" validate:"min=0,max=100"`
	ImageURL           *string     `gorm:"size:255" json:"image_url"`
	FormConfig         []FormField `gorm:"type:jsonb;serializer:json;not null" json:"form_config" validate:"required,dive"`
	Status             string      `gorm:"size:20;default:active;index:idx_category_status,idx_status_order" json:"status" validate:"oneof=active inactive out_of_stock"`
	DisplayOrder       int         `gorm:"default:0;index:idx_status_order" json:"display_order"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`

	// Relations
	Category     Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:ProductID" json:"transactions,omitempty"`
}

// GameAccount represents a user's game account
type GameAccount struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_game" json:"user_id" validate:"required"`
	GameName  string    `gorm:"size:100;not null;index:idx_user_game" json:"game_name" validate:"required,max=100"`
	GameID    string    `gorm:"size:100;not null" json:"game_id" validate:"required,max=100"`
	Server    *string   `gorm:"size:50" json:"server" validate:"omitempty,max=50"`
	ZoneID    *string   `gorm:"size:50" json:"zone_id" validate:"omitempty,max=50"`
	Nickname  *string   `gorm:"size:100" json:"nickname" validate:"omitempty,max=100"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook for GameAccount
func (ga *GameAccount) BeforeCreate(tx *gorm.DB) error {
	// If this is set as primary, unset other primary accounts for the same user and game
	if ga.IsPrimary {
		tx.Model(&GameAccount{}).
			Where("user_id = ? AND game_name = ? AND is_primary = ?", ga.UserID, ga.GameName, true).
			Update("is_primary", false)
	}
	return nil
}

// BeforeUpdate hook for GameAccount
func (ga *GameAccount) BeforeUpdate(tx *gorm.DB) error {
	// If this is set as primary, unset other primary accounts for the same user and game
	if ga.IsPrimary {
		tx.Model(&GameAccount{}).
			Where("user_id = ? AND game_name = ? AND id != ? AND is_primary = ?", ga.UserID, ga.GameName, ga.ID, true).
			Update("is_primary", false)
	}
	return nil
}
