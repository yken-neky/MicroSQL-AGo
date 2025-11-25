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
	// attach request-scoped logging middleware
	logMW := middleware.NewLoggingMiddleware(logger)
	r.Use(logMW.RequestLogger())
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

	// Auth routes (login unprotected, logout protected)
	auth := api.Group("/auth")
	{
		uh := handlers.NewUserHandlerWithJWT(db, logger, jwtService)
		auth.POST("/login", uh.Login)
		// logout is protected: user must include valid bearer token
		authMW := middleware.NewAuthMiddleware(jwtService)
		auth.POST("/logout", authMW.RequireAuth(), uh.Logout)
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

		// add audit repository for persistence of runs/results
		auditRepo := repo.NewGormAuditRepository(db)
		auditUC := controlsuc.NewExecuteAuditUseCase(controlsRepo, sqlService, queryExec, connRepo, auditRepo)
		ah := handlers.NewAuditHandler(auditUC)

		audits.POST("/execute", ah.ExecuteAudit)
		// fetch audit run details
		audits.GET("/:id", ah.GetAudit)
	}

	// Admin endpoints
	admin := api.Group("/admin")
	{
		adminMW := middleware.NewAuthMiddleware(jwtService)
		// only admin role may access admin routes
		admin.Use(adminMW.RequireRole("admin"))
		sessionRepo := repo.NewGormSessionRepository(db)
		adminHandler := handlers.NewAdminHandler(db, logger, sessionRepo)
		admin.GET("/sessions", adminHandler.ListActiveSessions)
	}

	// stub connections endpoint to avoid 404s (real implementation lives elsewhere)
	api.GET("/connections", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "connections endpoint not implemented yet"})
	})
}
