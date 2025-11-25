package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"gorm.io/gorm"
)

// AdminHandler for admin-only endpoints
type AdminHandler struct {
	DB          *gorm.DB
	Logger      *zap.Logger
	SessionRepo repositories.SessionRepository
}

func NewAdminHandler(db *gorm.DB, logger *zap.Logger, sr repositories.SessionRepository) *AdminHandler {
	return &AdminHandler{DB: db, Logger: logger, SessionRepo: sr}
}

// ListActiveSessions returns all active sessions with user info (admin only)
func (h *AdminHandler) ListActiveSessions(c *gin.Context) {
	if h.SessionRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session repository not configured"})
		return
	}

	sessions, err := h.SessionRepo.ListActiveSessions()
	if err != nil {
		h.Logger.Error("failed listing active sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	resp := make([]map[string]interface{}, 0, len(sessions))
	for _, s := range sessions {
		// fetch user info
		var u entities.User
		if err := h.DB.First(&u, s.UserID).Error; err != nil {
			// if user not found, still include session, with null user
			resp = append(resp, map[string]interface{}{
				"session_id": s.ID,
				"user_id":    s.UserID,
				"username":   nil,
				"email":      nil,
				"token":      s.Token,
				"expires_at": s.ExpiresAt,
				"created_at": s.CreatedAt,
			})
			continue
		}
		resp = append(resp, map[string]interface{}{
			"session_id": s.ID,
			"user_id":    s.UserID,
			"username":   u.Username,
			"email":      u.Email,
			"token":      s.Token,
			"expires_at": s.ExpiresAt,
			"created_at": s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": resp})
}
