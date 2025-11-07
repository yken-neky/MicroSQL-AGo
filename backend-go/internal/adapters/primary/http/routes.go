package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterRoutes wires HTTP routes. Keep minimal for now.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// register user routes
	api := r.Group("/api")
	users := api.Group("/users")
	{
		// handlers wired later; placeholder health route kept
		users.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	}
}
