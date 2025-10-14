package middlewares

import (
	"time"

	"temp/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware logs HTTP requests using zap logger
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log with different levels based on status code
		if status >= 500 {
			global.Logger.Error("HTTP Request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		} else if status >= 400 {
			global.Logger.Warn("HTTP Request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		} else {
			global.Logger.Info("HTTP Request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
			)
		}
	}
}
