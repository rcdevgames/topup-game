package services

import (
	"gorm.io/gorm"

	"topup-backend/internal/config"
)

// ServiceContainer holds all service dependencies
type ServiceContainer struct {
	DB     *gorm.DB
	Config *config.Config

	// Services
	UserService        *UserService
	AdminService       *AdminService
	CategoryService    *CategoryService
	ProductService     *ProductService
	GameAccountService *GameAccountService
	VoucherService     *VoucherService
	TransactionService *TransactionService
	FileUploadService  *FileUploadService
	CacheService       *CacheService
	WhatsAppService    *WhatsAppService
	PaymentService     *PaymentService
	AnalyticsService   *AnalyticsService
}

// NewServiceContainer creates a new service container
func NewServiceContainer(db *gorm.DB, cfg *config.Config) *ServiceContainer {
	// Initialize cache service
	cacheService := NewCacheService(cfg)

	// Initialize file upload service
	fileUploadService := NewFileUploadService(cfg)

	// Initialize WhatsApp service
	whatsappService := NewWhatsAppService(cfg)

	// Initialize payment service
	paymentService := NewPaymentService(cfg)

	// Initialize services
	userService := NewUserService(db, cfg)
	adminService := NewAdminService(db, cfg)
	categoryService := NewCategoryService(db, cacheService)
	productService := NewProductService(db, cacheService)
	gameAccountService := NewGameAccountService(db)
	voucherService := NewVoucherService(db, cacheService)
	transactionService := NewTransactionService(db, voucherService, paymentService, whatsappService, cfg)
	analyticsService := NewAnalyticsService(db, cacheService)

	return &ServiceContainer{
		DB:     db,
		Config: cfg,

		UserService:        userService,
		AdminService:       adminService,
		CategoryService:    categoryService,
		ProductService:     productService,
		GameAccountService: gameAccountService,
		VoucherService:     voucherService,
		TransactionService: transactionService,
		FileUploadService:  fileUploadService,
		CacheService:       cacheService,
		WhatsAppService:    whatsappService,
		PaymentService:     paymentService,
		AnalyticsService:   analyticsService,
	}
}
