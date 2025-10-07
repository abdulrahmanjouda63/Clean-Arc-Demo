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
	"strings"

	"temp/config"
	"temp/docs"
	"temp/global"
	"temp/handlers"
	"temp/repositories"
	"temp/routes"
	"temp/services"
)

func main() {
	log.Println("Starting application initialization...")

	// Load configuration first to know port
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Println("Configuration loaded.")

	// Resolve a free port: if configured port is busy, pick a free one
	desiredPort := cfg.Server.Port
	if !strings.HasPrefix(desiredPort, ":") {
		desiredPort = ":" + desiredPort
	}

	// Try the desired port first
	finalPort := desiredPort
	if ln, errListen := net.Listen("tcp", finalPort); errListen == nil {
		_ = ln.Close()
	} else {
		// Fallback to an ephemeral free port
		if ln2, errListen2 := net.Listen("tcp", "127.0.0.1:0"); errListen2 == nil {
			addr := ln2.Addr().(*net.TCPAddr)
			finalPort = fmt.Sprintf(":%d", addr.Port)
			_ = ln2.Close()
			log.Printf("Configured port %s is busy; falling back to free port %s", desiredPort, finalPort)
		} else {
			log.Fatalf("unable to acquire a free port: %v", errListen2)
		}
	}
	cfg.Server.Port = finalPort

	// Initialize Swagger docs with the resolved port
	docs.SwaggerInfo.Title = "User Management API"
	docs.SwaggerInfo.Description = "A comprehensive API for user management with authentication, registration, and profile management"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = "localhost" + cfg.Server.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// config already loaded above

	// init db
	log.Println("Initializing database...")
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	global.DB = db
	log.Println("Database initialized.")

	// Init Redis
	log.Println("Initializing Redis...")
	if err = config.InitRedis(cfg); err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}
	// global.RedisClient = rdb // No longer needed as InitRedis directly sets global.Redis
	log.Println("Redis initialized.")

	// repos, services, handlers
	log.Println("Initializing repositories, services, and handlers...")
	userRepo := repositories.NewUserRepo()
	if err := userRepo.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	userService := services.NewUserService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	userHandler := handlers.NewUserHandler(userService)
	log.Println("Repositories, services, and handlers initialized.")

	r := routes.NewRouter(userHandler, cfg.JWT.Secret)

	fmt.Println("Server listening on", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
