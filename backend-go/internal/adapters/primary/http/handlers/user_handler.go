package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	dto "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/dto"
)

// UserHandler handles user related endpoints
type UserHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{DB: db, Logger: logger}
}

// Register registers a new user. Minimal wiring; delegates to usecase in next steps.
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: call usecase
	c.JSON(http.StatusCreated, gin.H{"message": "user registered (stub)"})
}

// Login logs the user in (stub)
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: call usecase
	c.JSON(http.StatusOK, gin.H{"token": "stub-token"})
}
