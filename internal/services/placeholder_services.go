package services

import (
	"topup-backend/internal/config"

	"gorm.io/gorm"
)

// GameAccountService handles game account-related business logic
type GameAccountService struct {
	db *gorm.DB
}

// NewGameAccountService creates a new game account service
func NewGameAccountService(db *gorm.DB) *GameAccountService {
	return &GameAccountService{
		db: db,
	}
}

// VoucherService handles voucher-related business logic
type VoucherService struct {
	db    *gorm.DB
	cache *CacheService
}

// NewVoucherService creates a new voucher service
func NewVoucherService(db *gorm.DB, cache *CacheService) *VoucherService {
	return &VoucherService{
		db:    db,
		cache: cache,
	}
}

// TransactionService handles transaction-related business logic
type TransactionService struct {
	db              *gorm.DB
	voucherService  *VoucherService
	paymentService  *PaymentService
	whatsappService *WhatsAppService
	cfg             *config.Config
}

// NewTransactionService creates a new transaction service
func NewTransactionService(db *gorm.DB, voucherService *VoucherService, paymentService *PaymentService, whatsappService *WhatsAppService, cfg *config.Config) *TransactionService {
	return &TransactionService{
		db:              db,
		voucherService:  voucherService,
		paymentService:  paymentService,
		whatsappService: whatsappService,
		cfg:             cfg,
	}
}

// WhatsAppService handles WhatsApp integration
type WhatsAppService struct {
	cfg *config.Config
}

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(cfg *config.Config) *WhatsAppService {
	return &WhatsAppService{
		cfg: cfg,
	}
}

// PaymentService handles payment gateway integration
type PaymentService struct {
	cfg *config.Config
}

// NewPaymentService creates a new payment service
func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{
		cfg: cfg,
	}
}

// AnalyticsService handles analytics and reporting
type AnalyticsService struct {
	db    *gorm.DB
	cache *CacheService
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *gorm.DB, cache *CacheService) *AnalyticsService {
	return &AnalyticsService{
		db:    db,
		cache: cache,
	}
}
