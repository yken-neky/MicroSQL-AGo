package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	dto "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/dto"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/security"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user related endpoints
type UserHandler struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	JWTService *security.JWTService
}

func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{DB: db, Logger: logger}
}

// NewUserHandlerWithJWT creates UserHandler with JWT service
func NewUserHandlerWithJWT(db *gorm.DB, logger *zap.Logger, jwtService *security.JWTService) *UserHandler {
	return &UserHandler{
		DB:         db,
		Logger:     logger,
		JWTService: jwtService,
	}
}

// Register registers a new user. Minimal wiring; delegates to usecase in next steps.
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check username uniqueness
	var existing entities.User
	if err := h.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		h.Logger.Error("database error checking username", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Check email uniqueness
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		h.Logger.Error("database error checking email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := entities.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Ensure role has sensible default if not set by DB defaults
	if user.Role == "" {
		user.Role = "cliente"
	}

	// Omit LastLogin on insert to avoid zero-datetime errors on strict MySQL modes.
	if err := h.DB.Omit("last_login").Create(&user).Error; err != nil {
		h.Logger.Error("failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generate JWT token for the new user
	if h.JWTService == nil {
		h.Logger.Warn("JWTService is nil, cannot generate token for new user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication service not available"})
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		h.Logger.Error("failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Update LastLogin to now
	h.DB.Model(&user).Update("last_login", gorm.Expr("CURRENT_TIMESTAMP"))

	userResp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	c.JSON(http.StatusCreated, dto.LoginResponse{Token: token, User: userResp})
}

// Login logs the user in and returns JWT token
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user entities.User
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.Logger.Error("database error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user account is inactive"})
		return
	}

	// Generate JWT token
	if h.JWTService == nil {
		h.Logger.Warn("JWTService is nil, cannot generate token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication service not available"})
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		h.Logger.Error("failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Update LastLogin
	h.DB.Model(&user).Update("last_login", gorm.Expr("CURRENT_TIMESTAMP"))

	// Return response
	userResp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  userResp,
	})
}
