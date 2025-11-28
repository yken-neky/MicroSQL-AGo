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
	connectionuc "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/connection"
	controlsuc "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/controls"
	repo "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	sqlexec "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/sqlserver"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/encryption"
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

	// NOTE: Audit endpoints are attached under /api/db/:manager/audits to make the manager explicit

	// Admin endpoints
	admin := api.Group("/admin")
	{
		adminMW := middleware.NewAuthMiddleware(jwtService)
		// only admin role may access admin routes
		admin.Use(adminMW.RequireRole("admin"))
		sessionRepo := repo.NewGormSessionRepository(db)
		adminHandler := handlers.NewAdminHandler(db, logger, sessionRepo)
		admin.GET("/sessions", adminHandler.ListActiveSessions)
		// metrics endpoints for admin
		admin.GET("/metrics/users", adminHandler.GetUsersMetrics)
		admin.GET("/metrics/connections", adminHandler.GetConnectionsMetrics)
		admin.GET("/metrics/audits", adminHandler.GetAuditsMetrics)
		admin.GET("/metrics/roles", adminHandler.GetRolesMetrics)
		admin.GET("/metrics/system", adminHandler.GetSystemMetrics)
	}

	// DB connection endpoints: /api/db and /api/db/:manager
	dbGroup := api.Group("/db")
	{
		connAuth := middleware.NewAuthMiddleware(jwtService)
		dbGroup.Use(connAuth.RequireAuth())

		connRepo := repo.NewGormConnectionRepository(db)
		sqlService := sqladp.NewSQLServerAdapter(logger)
		// encryption service for persisting DB passwords
		encService := encryption.NewAESGCMService(cfg.EncKey)

		connectUC := connectionuc.NewConnectToServerUseCase(connRepo, sqlService, encService)
		disconnectUC := connectionuc.NewDisconnectFromServerUseCase(connRepo, sqlService)
		getActiveUC := connectionuc.NewGetActiveConnectionUseCase(connRepo)
		listUC := connectionuc.NewListActiveConnectionsUseCase(connRepo)

		// history UC to list connection logs
		historyUC := connectionuc.NewListConnectionHistoryUseCase(repo.NewGormConnectionRepository(db))
		ch := handlers.NewConnectionHandler(connectUC, disconnectUC, getActiveUC, listUC, historyUC, logger)

		// List all active connections for user across drivers
		dbGroup.GET("/connections", ch.GetActive)

		// per-manager operations
		mgr := dbGroup.Group(":manager")
		{
			// open a connection: POST /api/db/:manager/open
			mgr.POST("/open", ch.Connect)
			// close connection for manager: DELETE /api/db/:manager/close
			mgr.DELETE("/close", ch.Disconnect)
			// get active connection for manager: GET /api/db/:manager/connection
			mgr.GET("/connection", ch.GetActive)

			// Audits routes under the explicit manager: /api/db/:manager/audits
			controlsRepo := repo.NewGormControlsRepository(db)
			queryExec := sqlexec.NewSQLServerQueryExecutor()
			auditRepo := repo.NewGormAuditRepository(db)
			auditUC := controlsuc.NewExecuteAuditUseCase(controlsRepo, sqlService, queryExec, connRepo, auditRepo, encService)
			ah := handlers.NewAuditHandler(auditUC)

			mgr.POST("/audits/execute", ah.ExecuteAudit)
			mgr.GET("/audits/:id", ah.GetAudit)
		}
	}
}
