package routes

import (
	"temp/handlers"
	"temp/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userHandler *handlers.UserHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	// Restrict trusted proxies to loopback by default to avoid security warning
	// and prevent trusting all proxies.
	// This allows running safely in dev on localhost. Adjust for production as needed.
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1", "localhost"})
	r.Use(gin.Recovery())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.RecoveryMiddleware())

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.GET("/profile", middlewares.AuthMiddleware(jwtSecret), userHandler.Profile)
		v1.POST("/set-redis-key", userHandler.SetRedisKey)
	}

	return r
}
