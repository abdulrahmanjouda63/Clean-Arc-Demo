// Package main provides a REST API for user management
// @title User Management API
// @version 1.0.0
// @description A comprehensive API for user management with authentication, registration, and profile management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"temp/cli"
	"temp/config"
	"temp/docs"
	"temp/global"
	"temp/handlers"
	"temp/repositories"
	"temp/routes"
	"temp/services"

	"go.uber.org/zap"
)

func main() {
	// Load configuration first
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize zap logger with configuration
	if err := global.InitLogger(cfg.Logging.Level, cfg.Logging.OutputPath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer global.SyncLogger()
	
	global.Logger.Info("Starting application initialization...")

	// Handle CLI commands
	if handled, err := cli.ParseAndExecute(cfg); handled {
		if err != nil {
			global.Logger.Fatal("CLI command failed", zap.Error(err))
		}
		return
	}

	// Resolve a free port: if configured port is busy, pick a free one
	finalPort := resolvePort(cfg.Server.Port)
	cfg.Server.Port = finalPort

	// Initialize Swagger docs
	initSwagger(cfg.Server.Port)

	// Initialize database
	global.Logger.Info("Initializing database...")
	db, err := config.InitDB(cfg)
	if err != nil {
		global.Logger.Fatal("DB init error", zap.Error(err))
	}
	global.DB = db
	global.Logger.Info("Database initialized.")

	// Initialize Redis (with Sentinel support)
	global.Logger.Info("Initializing Redis...")
	if err = config.InitRedis(cfg); err != nil {
		global.Logger.Fatal("Error initializing Redis", zap.Error(err))
	}
	defer config.CloseRedis()
	global.Logger.Info("Redis initialized.")

	// Initialize repositories
	global.Logger.Info("Initializing repositories, services, and handlers...")
	userRepo := repositories.NewUserRepo()
	if err := userRepo.Migrate(); err != nil {
		global.Logger.Fatal("Migration failed", zap.Error(err))
	}

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	emailService := services.NewEmailService(
		cfg.Email.Enabled,
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.From,
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	global.Logger.Info("Repositories, services, and handlers initialized.")

	// Create router with CORS configuration
	r := routes.NewRouter(userHandler, cfg.JWT.Secret, cfg)

	// Setup graceful shutdown
	setupGracefulShutdown()

	// Start server
	fmt.Printf("Server listening on %s\n", cfg.Server.Port)
	fmt.Printf("Swagger UI available at http://localhost%s/swagger/index.html\n", cfg.Server.Port)
	global.Logger.Info("Server listening", zap.String("port", cfg.Server.Port))
	
	if err := r.Run(cfg.Server.Port); err != nil {
		global.Logger.Fatal("Failed to run server", zap.Error(err))
	}

	// Send test welcome email (optional)
	if cfg.Email.Enabled {
		go func() {
			// Example: Send welcome email to new users
			_ = emailService.SendWelcomeEmail("test@example.com", "Test User")
		}()
	}
}

// resolvePort resolves a free port if the configured one is busy
func resolvePort(desiredPort string) string {
	if !strings.HasPrefix(desiredPort, ":") {
		desiredPort = ":" + desiredPort
	}

	// Try the desired port first
	if ln, err := net.Listen("tcp", desiredPort); err == nil {
		_ = ln.Close()
		return desiredPort
	}

	// Fallback to an ephemeral free port
	if ln, err := net.Listen("tcp", "127.0.0.1:0"); err == nil {
		addr := ln.Addr().(*net.TCPAddr)
		finalPort := fmt.Sprintf(":%d", addr.Port)
		_ = ln.Close()
		global.Logger.Warn("Configured port is busy; falling back to free port",
			zap.String("desired", desiredPort),
			zap.String("final", finalPort),
		)
		return finalPort
	}

	global.Logger.Fatal("Unable to acquire a free port")
	return desiredPort
}

// initSwagger initializes Swagger documentation
func initSwagger(port string) {
	docs.SwaggerInfo.Title = "User Management API"
	docs.SwaggerInfo.Description = "A comprehensive API for user management with authentication, registration, and profile management"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = "localhost" + port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

// setupGracefulShutdown handles graceful shutdown on interrupt signals
func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		global.Logger.Info("Shutting down gracefully...")
		
		// Close Redis connection
		if err := config.CloseRedis(); err != nil {
			global.Logger.Error("Error closing Redis", zap.Error(err))
		}
		
		// Close database connection
		if global.DB != nil {
			if sqlDB, err := global.DB.DB(); err == nil {
				if err := sqlDB.Close(); err != nil {
					global.Logger.Error("Error closing database", zap.Error(err))
				}
			}
		}
		
		global.Logger.Info("Shutdown complete")
		os.Exit(0)
	}()
}