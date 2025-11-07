package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	httpadp "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/primary/http"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/persistence/sqlite/migrations"
	cfg "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/config"
	pkglog "github.com/yken-neky/MicroSQL-AGo/backend-go/pkg/utils"
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

	r := gin.Default()
	httpadp.RegisterRoutes(r, db, logger)

	addr := fmt.Sprintf(":%s", cfgVal.ServerPort)
	logger.Info("starting server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited with: %v", err)
	}
}
