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
    "strings"

    "temp/config"
    "temp/docs"
    "temp/global"
    "temp/handlers"
    "temp/repositories"
    "temp/routes"
    "temp/services"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func main() {
    // Initialize zap logger
    if err := global.InitLogger(); err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer global.Logger.Sync()
    global.Logger.Info("Starting application initialization...")

    // Load configuration first to know port
    global.Logger.Info("Loading configuration...")
    cfg, err := config.LoadConfig("./config.yaml")
    if err != nil {
        global.Logger.Fatal("Error loading configuration", zap.Error(err))
    }
    global.Logger.Info("Configuration loaded.")

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
            global.Logger.Warn("Configured port is busy; falling back to free port", zap.String("desired", desiredPort), zap.String("final", finalPort))
        } else {
            global.Logger.Fatal("Unable to acquire a free port", zap.Error(errListen2))
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

    // init db
    global.Logger.Info("Initializing database...")
    db, err := config.InitDB(cfg)
    if err != nil {
        global.Logger.Fatal("DB init error", zap.Error(err))
    }
    global.DB = db
    global.Logger.Info("Database initialized.")

    // Init Redis
    global.Logger.Info("Initializing Redis...")
    if err = config.InitRedis(cfg); err != nil {
        global.Logger.Fatal("Error initializing Redis", zap.Error(err))
    }
    global.Logger.Info("Redis initialized.")

    // repos, services, handlers
    global.Logger.Info("Initializing repositories, services, and handlers...")
    userRepo := repositories.NewUserRepo()
    if err := userRepo.Migrate(); err != nil {
        global.Logger.Fatal("Migration failed", zap.Error(err))
    }

    userService := services.NewUserService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
    userHandler := handlers.NewUserHandler(userService)
    global.Logger.Info("Repositories, services, and handlers initialized.")

    // Gin router with CORS
    r := gin.Default()
    r.Use(cors.Default())

    // Register routes
    routes.RegisterRoutes(r, userHandler, cfg.JWT.Secret)

    // CLI support (migrate, etc.)
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "migrate":
            if err := userRepo.Migrate(); err != nil {
                global.Logger.Fatal("Migration failed", zap.Error(err))
            }
            global.Logger.Info("Migration completed successfully.")
            return
        }
    }

    fmt.Println("Server listening on", cfg.Server.Port)
    global.Logger.Info("Server listening", zap.String("port", cfg.Server.Port))
    if err := r.Run(cfg.Server.Port); err != nil {
        global.Logger.Fatal("Failed to run server", zap.Error(err))
    }
}