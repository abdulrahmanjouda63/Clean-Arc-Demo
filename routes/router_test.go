package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func TestSwaggerEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a simple router to test Swagger endpoint
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Test Swagger endpoint
	req := httptest.NewRequest("GET", "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	// Should return 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}
