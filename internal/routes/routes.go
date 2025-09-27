package routes

import (
	"github.com/gin-gonic/gin"

	"topup-backend/internal/config"
	"topup-backend/internal/handlers"
	"topup-backend/internal/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, handlers *handlers.HandlerContainer, cfg *config.Config) {
	// API v1 group
	api := router.Group("/api")
	{
		// Public routes (no authentication required)
		setupPublicRoutes(api, handlers, cfg)

		// Protected routes (authentication required)
		setupProtectedRoutes(api, handlers, cfg)

		// Admin routes (admin authentication required)
		setupAdminRoutes(api, handlers, cfg)
	}
}

// setupPublicRoutes sets up public routes that don't require authentication
func setupPublicRoutes(api *gin.RouterGroup, handlers *handlers.HandlerContainer, cfg *config.Config) {
	// CSRF routes (must be accessible without authentication)
	csrf := api.Group("/csrf")
	{
		csrf.GET("", handlers.CSRFHandler.GetCSRFToken)
		csrf.POST("/validate", middleware.CSRFCheck(), handlers.CSRFHandler.ValidateCSRF)
	}

	// Authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.UserHandler.Register)
		auth.POST("/login", handlers.UserHandler.Login)
		auth.POST("/refresh-token", handlers.UserHandler.RefreshToken)
	}

	// Public product and category routes
	categories := api.Group("/categories")
	{
		categories.GET("", handlers.CategoryHandler.GetAllCategories)
		categories.GET("/:id", handlers.CategoryHandler.GetCategoryByID)
	}

	products := api.Group("/products")
	{
		products.GET("", handlers.ProductHandler.GetAllProducts)
		products.GET("/:id", handlers.ProductHandler.GetProductByID)
		products.GET("/search", handlers.ProductHandler.SearchProducts)
		products.GET("/category/:category_id", handlers.ProductHandler.GetProductsByCategory)
	}

	// Public voucher validation
	vouchers := api.Group("/vouchers")
	{
		vouchers.POST("/validate", handlers.VoucherHandler.ValidateVoucher)
	}
}

// setupProtectedRoutes sets up routes that require user authentication
func setupProtectedRoutes(api *gin.RouterGroup, handlers *handlers.HandlerContainer, cfg *config.Config) {
	// Apply authentication middleware and CSRF protection
	protected := api.Group("", middleware.AuthMiddleware(cfg), middleware.CSRFCheck())

	// User profile routes
	auth := protected.Group("/auth")
	{
		auth.POST("/logout", handlers.UserHandler.Logout)
		auth.GET("/profile", handlers.UserHandler.GetProfile)
		auth.PUT("/profile", handlers.UserHandler.UpdateProfile)
		auth.POST("/change-password", handlers.UserHandler.ChangePassword)
	}

	// Game accounts management
	gameAccounts := protected.Group("/game-accounts")
	{
		gameAccounts.GET("", handlers.GameAccountHandler.GetUserGameAccounts)
		gameAccounts.POST("", handlers.GameAccountHandler.CreateGameAccount)
		gameAccounts.PUT("/:id", handlers.GameAccountHandler.UpdateGameAccount)
		gameAccounts.DELETE("/:id", handlers.GameAccountHandler.DeleteGameAccount)
	}

	// Transaction routes
	transactions := protected.Group("/transactions")
	{
		transactions.POST("", handlers.TransactionHandler.CreateTransaction)
		transactions.GET("", handlers.TransactionHandler.GetUserTransactions)
		transactions.GET("/:id", handlers.TransactionHandler.GetTransactionByID)
		transactions.POST("/:id/cancel", handlers.TransactionHandler.CancelTransaction)
	}
}

// setupAdminRoutes sets up routes that require admin authentication
func setupAdminRoutes(api *gin.RouterGroup, handlers *handlers.HandlerContainer, cfg *config.Config) {
	// Admin authentication routes (no middleware)
	adminAuth := api.Group("/admin/auth")
	{
		adminAuth.POST("/login", handlers.AdminHandler.Login)
		adminAuth.POST("/refresh-token", handlers.AdminHandler.RefreshToken)
	}

	// Protected admin routes with CSRF protection
	admin := api.Group("/admin", middleware.AdminAuthMiddleware(cfg), middleware.CSRFCheck())
	{
		// Admin profile
		admin.POST("/auth/logout", handlers.AdminHandler.Logout)
		admin.GET("/auth/profile", handlers.AdminHandler.GetProfile)

		// Dashboard analytics
		admin.GET("/analytics/dashboard", handlers.AnalyticsHandler.GetDashboardAnalytics)

		// Admin users management (super_admin and admin only)
		adminUsers := admin.Group("/users", middleware.RoleMiddleware("super_admin", "admin"))
		{
			adminUsers.GET("", handlers.AdminHandler.GetAllAdmins)
			adminUsers.POST("", handlers.AdminHandler.CreateAdmin)
			adminUsers.GET("/:id", handlers.AdminHandler.GetAdminByID)
			adminUsers.PUT("/:id", handlers.AdminHandler.UpdateAdmin)
			adminUsers.DELETE("/:id", handlers.AdminHandler.DeleteAdmin)
		}

		// Categories management
		categories := admin.Group("/categories")
		{
			categories.GET("", handlers.CategoryHandler.GetCategoriesForAdmin)
			categories.POST("", middleware.RoleMiddleware("super_admin", "admin"), handlers.CategoryHandler.CreateCategory)
			categories.PUT("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.CategoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.CategoryHandler.DeleteCategory)
			categories.POST("/:id/upload-icon", middleware.RoleMiddleware("super_admin", "admin"), handlers.FileUploadHandler.UploadCategoryIcon)
		}

		// Products management
		products := admin.Group("/products")
		{
			products.GET("", handlers.ProductHandler.GetProductsForAdmin)
			products.POST("", middleware.RoleMiddleware("super_admin", "admin"), handlers.ProductHandler.CreateProduct)
			products.GET("/:id", handlers.ProductHandler.GetProductByIDForAdmin)
			products.PUT("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.ProductHandler.UpdateProduct)
			products.DELETE("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.ProductHandler.DeleteProduct)
			products.POST("/:id/upload-image", middleware.RoleMiddleware("super_admin", "admin"), handlers.FileUploadHandler.UploadProductImage)
		}

		// Vouchers management
		vouchers := admin.Group("/vouchers")
		{
			vouchers.GET("", handlers.VoucherHandler.GetVouchersForAdmin)
			vouchers.POST("", middleware.RoleMiddleware("super_admin", "admin"), handlers.VoucherHandler.CreateVoucher)
			vouchers.GET("/:id", handlers.VoucherHandler.GetVoucherByID)
			vouchers.PUT("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.VoucherHandler.UpdateVoucher)
			vouchers.DELETE("/:id", middleware.RoleMiddleware("super_admin", "admin"), handlers.VoucherHandler.DeleteVoucher)
			vouchers.GET("/:id/usage-stats", handlers.VoucherHandler.GetVoucherUsageStats)
		}

		// Transactions management
		transactions := admin.Group("/transactions")
		{
			transactions.GET("", handlers.TransactionHandler.GetTransactionsForAdmin)
			transactions.GET("/:id", handlers.TransactionHandler.GetTransactionByIDForAdmin)
			transactions.PUT("/:id/status", handlers.TransactionHandler.UpdateTransactionStatus)
			transactions.GET("/export", handlers.TransactionHandler.ExportTransactions)
		}

		// Customer users management
		customers := admin.Group("/customers")
		{
			customers.GET("", handlers.UserHandler.GetUserByID) // Reuse for now
		}
	}
}
