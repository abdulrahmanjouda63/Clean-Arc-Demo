package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"temp/global"
	"temp/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		global.Logger.Error("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			global.Logger.Error("Registration failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	global.Logger.Info("User registered successfully", 
		zap.Uint("user_id", u.ID),
		zap.String("email", u.Email),
	)

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
		global.Logger.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, u, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		global.Logger.Warn("Login failed", 
			zap.String("email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	global.Logger.Info("User logged in successfully",
		zap.Uint("user_id", u.ID),
		zap.String("email", u.Email),
	)

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
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	global.Logger.Info("Profile accessed", zap.Uint("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "this is a protected route", "user_id": float64(userID)})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string} true "Profile update data"
// @Success 200 {object} object{message=string} "Profile updated successfully"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 401 {object} object{error=string} "Unauthorized - invalid or missing token"
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	global.Logger.Info("Profile updated", zap.Uint("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{old_password=string,new_password=string} true "Password change data"
// @Success 200 {object} object{message=string} "Password changed successfully"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 401 {object} object{error=string} "Unauthorized - invalid or missing token"
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	global.Logger.Info("Password changed", zap.Uint("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
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
		global.Logger.Error("Failed to set Redis key", 
			zap.String("key", req.Key),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set key in redis: %v", err)})
		return
	}

	global.Logger.Info("Redis key set", zap.String("key", req.Key))
	c.JSON(http.StatusOK, gin.H{"message": "key set successfully", "key": req.Key, "value": req.Value})
}

// GetRedisKey godoc
// @Summary Get Redis value by key
// @Description Retrieve a value from Redis cache by key
// @Tags Redis
// @Accept json
// @Produce json
// @Param key path string true "Redis key"
// @Success 200 {object} object{key=string,value=string} "Value retrieved successfully"
// @Failure 404 {object} object{error=string} "Key not found"
// @Failure 500 {object} object{error=string} "Internal server error - Redis operation failed"
// @Router /get-redis-key/{key} [get]
func (h *UserHandler) GetRedisKey(c *gin.Context) {
	key := c.Param("key")
	
	ctx := c.Request.Context()
	value, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		global.Logger.Error("Failed to get Redis key",
			zap.String("key", key),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get key from redis: %v", err)})
		return
	}

	global.Logger.Info("Redis key retrieved", zap.String("key", key))
	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

// Helper function to extract user ID from context
func getUserID(c *gin.Context) uint {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	switch v := userIDRaw.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		return uint(parsed)
	default:
		return 0
	}
}