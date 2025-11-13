package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	handlers "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/handlers"
	middleware "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http/middleware"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/security"
	sqladp "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/sqlserver"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/config"
	controlsuc "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/controls"
	repo "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	sqlexec "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/sqlserver"
)

// RegisterRoutes wires HTTP routes with JWT middleware.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	// Create JWT service from config
	cfg := config.LoadConfig()
	jwtService := security.NewJWTService(cfg.JWTSecret)

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

	// Auth routes (unprotected)
	auth := api.Group("/auth")
	{
		uh := handlers.NewUserHandlerWithJWT(db, logger, jwtService)
		auth.POST("/login", uh.Login)
	}

	// User routes
	users := api.Group("/users")
	{
		users.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
		uh := handlers.NewUserHandlerWithJWT(db, logger, jwtService)
		users.POST("/register", uh.Register)
	}

	// Audit endpoints - execute predefined control scripts (requires auth)
	audits := api.Group("/audits")
	{
		// Apply JWT auth middleware
		authMW := middleware.NewAuthMiddleware(jwtService)
		audits.Use(authMW.RequireAuth())

		// wire dependencies for audit
		controlsRepo := repo.NewGormControlsRepository(db)
		connRepo := repo.NewGormConnectionRepository(db)

		sqlService := sqladp.NewSQLServerAdapter(logger)
		queryExec := sqlexec.NewSQLServerQueryExecutor()

		auditUC := controlsuc.NewExecuteAuditUseCase(controlsRepo, sqlService, queryExec, connRepo)
		ah := handlers.NewAuditHandler(auditUC)

		audits.POST("/execute", ah.ExecuteAudit)
	}

	// stub connections endpoint to avoid 404s (real implementation lives elsewhere)
	api.GET("/connections", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "connections endpoint not implemented yet"})
	})
}
