package handlers

import (
	"topup-backend/internal/services"
)

// HandlerContainer holds all handler dependencies
type HandlerContainer struct {
	Services *services.ServiceContainer

	// Handlers
	UserHandler        *UserHandler
	AdminHandler       *AdminHandler
	CategoryHandler    *CategoryHandler
	ProductHandler     *ProductHandler
	GameAccountHandler *GameAccountHandler
	VoucherHandler     *VoucherHandler
	TransactionHandler *TransactionHandler
	FileUploadHandler  *FileUploadHandler
	AnalyticsHandler   *AnalyticsHandler
	CSRFHandler        *CSRFHandler
}

// NewHandlerContainer creates a new handler container
func NewHandlerContainer(services *services.ServiceContainer) *HandlerContainer {
	return &HandlerContainer{
		Services: services,

		UserHandler:        NewUserHandler(services),
		AdminHandler:       NewAdminHandler(services),
		CategoryHandler:    NewCategoryHandler(services),
		ProductHandler:     NewProductHandler(services),
		GameAccountHandler: NewGameAccountHandler(services),
		VoucherHandler:     NewVoucherHandler(services),
		TransactionHandler: NewTransactionHandler(services),
		FileUploadHandler:  NewFileUploadHandler(services),
		AnalyticsHandler:   NewAnalyticsHandler(services),
		CSRFHandler:        NewCSRFHandler(),
	}
}
