package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "topup-backend/docs" // Import swagger docs
	"topup-backend/internal/config"
	"topup-backend/internal/database"
	"topup-backend/internal/handlers"
	"topup-backend/internal/middleware"
	"topup-backend/internal/routes"
	"topup-backend/internal/services"
)

// @title Waw Store Topup Game Online API
// @version 1.0
// @description Backend API for Waw Store - Game Topup Platform
// @termsOfService https://wawstore.com/terms

// @contact.name API Support
// @contact.url https://wawstore.com/support
// @contact.email support@wawstore.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	serviceContainer := services.NewServiceContainer(db, cfg)

	// Initialize handlers
	handlerContainer := handlers.NewHandlerContainer(serviceContainer)

	// Setup router
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	// Static file serving for local CDN
	router.Static("/cdn", cfg.Server.UploadDir)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Waw Store API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(router, handlerContainer, cfg)

	// Swagger documentation
	if cfg.Environment != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 Swagger docs available at: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
