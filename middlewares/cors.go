package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// CORSMiddleware returns a configured CORS middleware
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     config.AllowedMethods,
		AllowHeaders:     config.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}

// DefaultCORSMiddleware returns CORS middleware with default settings
func DefaultCORSMiddleware() gin.HandlerFunc {
	return CORSMiddleware(CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization", "Accept"},
	})
}
