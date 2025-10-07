package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		// simple print — يمكنك تبديله بلوغر متقدم
		println("[HTTP] ", method, path, "status:", status, "latency:", latency.String())
	}
}
