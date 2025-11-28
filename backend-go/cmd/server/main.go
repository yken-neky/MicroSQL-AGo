package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	httpadp "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/persistence/sqlite/migrations"
	cfg "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/config"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/persistence"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	pkglog "github.com/yken-neky/MicroSQL-AGo/backend-go/pkg/utils"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	cfgVal := cfg.LoadConfig()
	logger := pkglog.NewLogger(cfgVal.LogLevel)
	defer logger.Sync()

	db, err := cfg.NewGormDB(cfgVal)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	if err := migrations.Migrate(db); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// Seed initial roles & permissions (idempotent)
	if err := repositories.SeedDefaultRolesAndPermissions(db); err != nil {
		logger.Fatal("failed seeding roles/permissions", zap.Error(err))
	}

	// Attach a zap-backed GORM logger for structured SQL logging
	// Use a conservative slow query threshold (200ms) and Info level.
	gormLogger := persistence.NewZapGormLogger(logger, 200*time.Millisecond, gormlogger.Info)
	db = db.Session(&gorm.Session{Logger: gormLogger})

	r := gin.Default()
	httpadp.RegisterRoutes(r, db, logger)

	addr := fmt.Sprintf(":%s", cfgVal.ServerPort)
	logger.Info("starting server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited with: %v", err)
	}
}
