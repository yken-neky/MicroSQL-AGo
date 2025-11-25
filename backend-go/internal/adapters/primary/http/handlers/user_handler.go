package handlers

import (
	"net/http"
	"time"

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

	// Prevent login when user already has an active API session (single-session policy)
	var existingSession entities.Session
	if err := h.DB.Where("user_id = ? AND is_active = ?", user.ID, true).First(&existingSession).Error; err == nil {
		// if session exists and is not expired, disallow new login
		if existingSession.ExpiresAt == nil || existingSession.ExpiresAt.After(time.Now()) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already has an active session"})
			return
		}
		// otherwise mark as inactive and allow login
		existingSession.IsActive = false
		_ = h.DB.Save(&existingSession).Error
	} else if err != nil && err != gorm.ErrRecordNotFound {
		h.Logger.Error("failed checking existing session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
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

	// persist session for the generated token
	expires := time.Now().Add(24 * time.Hour)
	session := entities.Session{UserID: user.ID, Token: token, ExpiresAt: &expires, IsActive: true}
	if err := h.DB.Create(&session).Error; err != nil {
		h.Logger.Error("failed to persist session", zap.Error(err))
		// return internal error since session persistence is required for single-login enforcement
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

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

// Logout invalidates the currently presented token and marks the API session inactive.
// This endpoint must be called with Authorization: Bearer <token> and is protected by middleware.
func (h *UserHandler) Logout(c *gin.Context) {
	// Extract token from header
	authHeader := c.GetHeader("Authorization")
	tokenStr, err := security.ExtractBearerToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	// Get userID from context (set by auth middleware)
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return
	}

	// find active session for this token and user
	var s entities.Session
	if err := h.DB.Where("user_id = ? AND token = ? AND is_active = ?", userID, tokenStr, true).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// idempotent: nothing to do
			c.JSON(http.StatusOK, gin.H{"message": "no active session for token"})
			return
		}
		h.Logger.Error("failed reading session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	s.IsActive = false
	if err := h.DB.Save(&s).Error; err != nil {
		h.Logger.Error("failed deactivating session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
