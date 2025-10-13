package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"temp/services"

	"github.com/gin-gonic/gin"
	"temp/global"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(s services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: s}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{name=string,email=string,password=string} true "User registration data"
// @Success 201 {object} object{id=int,name=string,email=string} "User successfully registered"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 409 {object} object{error=string} "Conflict - user already exists"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "name": u.Name, "email": u.Email})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "User login credentials"
// @Success 200 {object} object{token=string,user=object{id=int,email=string,name=string}} "Login successful"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 401 {object} object{error=string} "Unauthorized - invalid credentials"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, u, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": u.ID, "email": u.Email, "name": u.Name}})
}

// Profile godoc
// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{message=string,user_id=int} "Profile retrieved successfully"
// @Failure 401 {object} object{error=string} "Unauthorized - invalid or missing token"
// @Router /profile [get]
func (h *UserHandler) Profile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	var userID uint
	switch v := userIDRaw.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case string:
		parsed, err := parseUintString(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
			return
		}
		userID = uint(parsed)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "this is a protected route", "user_id": float64(userID)})
}

// parseUintString tries to parse a string to uint64
func parseUintString(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// SetRedisKey godoc
// @Summary Set Redis key-value pair
// @Description Store a key-value pair in Redis cache
// @Tags Redis
// @Accept json
// @Produce json
// @Param request body object{key=string,value=string} true "Redis key-value data"
// @Success 200 {object} object{message=string,key=string,value=string} "Key set successfully"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 500 {object} object{error=string} "Internal server error - Redis operation failed"
// @Router /set-redis-key [post]
func (h *UserHandler) SetRedisKey(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := global.Redis.Set(ctx, req.Key, req.Value, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set key in redis: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key set successfully", "key": req.Key, "value": req.Value})
}
