package routes

import (
	"temp/handlers"
	"temp/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userHandler *handlers.UserHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1", "localhost"})

	r.Use(gin.Recovery())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.RecoveryMiddleware())
	r.Use(cors.Default())

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/profile", middlewares.AuthMiddleware(jwtSecret), userHandler.Profile)
		api.POST("/set-redis-key", userHandler.SetRedisKey)
		// Add more routes as needed
	}

	return r
}
