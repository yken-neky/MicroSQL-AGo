package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	handlers "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/handlers"
)

// RegisterRoutes wires HTTP routes. Keep minimal for now.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// root - help message
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "MicroSQL AGo backend", "status": "ok"})
	})

	// register user routes
	api := r.Group("/api")

	// swagger/info route
	api.GET("/swagger", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"swagger": "not generated"}) })

	users := api.Group("/users")
	{
		// minimal handlers: health, register, login
		users.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
		uh := handlers.NewUserHandler(db, logger)
		users.POST("/register", uh.Register)
		users.POST("/login", uh.Login)
	}

	// legacy/auth routes - map to same handlers to avoid 404s from clients
	auth := api.Group("/auth")
	{
		ah := handlers.NewUserHandler(db, logger)
		auth.POST("/login", ah.Login)
	}

	// stub connections endpoint to avoid 404s (real implementation lives elsewhere)
	api.GET("/connections", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "connections endpoint not implemented yet"})
	})
}
