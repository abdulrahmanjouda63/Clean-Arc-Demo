package routes

import (
	"temp/config"
	"temp/handlers"
	"temp/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userHandler *handlers.UserHandler, jwtSecret string, cfg *config.Config) *gin.Engine {
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1", "localhost"})

	// Recovery middleware
	r.Use(gin.Recovery())
	
	// Custom recovery middleware
	r.Use(middlewares.RecoveryMiddleware())
	
	// Logging middleware
	r.Use(middlewares.LoggingMiddleware())
	
	// CORS middleware
	if cfg.CORS.Enabled {
		r.Use(middlewares.CORSMiddleware(middlewares.CORSConfig{
			AllowedOrigins: cfg.CORS.AllowedOrigins,
			AllowedMethods: cfg.CORS.AllowedMethods,
			AllowedHeaders: cfg.CORS.AllowedHeaders,
		}))
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server is running",
		})
	})

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		
		// Protected routes
		protected := api.Group("")
		protected.Use(middlewares.AuthMiddleware(jwtSecret))
		{
			protected.GET("/profile", userHandler.Profile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.POST("/change-password", userHandler.ChangePassword)
		}
		
		// Redis operations (can be protected if needed)
		api.POST("/set-redis-key", userHandler.SetRedisKey)
		api.GET("/get-redis-key/:key", userHandler.GetRedisKey)
	}

	return r
}