package main

import (
    "flag"
    "fmt"
    "log"

    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/config"
    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// This migrator copies rows from a source SQLite DB into a destination MySQL DB.
// It accepts flags to configure source path and destination DSN (or uses env vars).
func main() {
    cfg := config.LoadConfig()

    srcPath := flag.String("src", cfg.DBPath, "path to source sqlite database file (overrides DB_PATH)")
    dstDSN := flag.String("dst-dsn", "", "MySQL DSN (user:pass@tcp(host:port)/dbname?params). If empty, uses MYSQL_* env vars from Config")
    flag.Parse()

    // Open source (sqlite)
    if *srcPath == "" {
        log.Fatalf("source sqlite path must be provided via --src or DB_PATH env")
    }
    srcDB, err := gorm.Open(sqlite.Open(*srcPath), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to open sqlite source DB '%s': %v", *srcPath, err)
    }

    // Build destination DSN
    var finalDSN string
    if *dstDSN != "" {
        finalDSN = *dstDSN
    } else {
        if cfg.MysqlHost == "" || cfg.MysqlUser == "" || cfg.MysqlDB == "" {
            log.Fatalf("either --dst-dsn or MYSQL_HOST, MYSQL_USER and MYSQL_DATABASE env vars must be set")
        }
        port := cfg.MysqlPort
        if port == "" {
            port = "3306"
        }
        finalDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            cfg.MysqlUser, cfg.MysqlPass, cfg.MysqlHost, port, cfg.MysqlDB,
        )
    }

    dstDB, err := gorm.Open(mysql.Open(finalDSN), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to open mysql destination DB: %v", err)
    }

    // Ensure destination has tables
    if err := dstDB.AutoMigrate(&entities.User{}, &entities.ActiveConnection{}, &entities.ControlsInformation{}); err != nil {
        log.Fatalf("failed to auto-migrate destination DB: %v", err)
    }

    // Migrate Users
    var users []entities.User
    if err := srcDB.Find(&users).Error; err != nil {
        log.Fatalf("failed to read users from sqlite: %v", err)
    }
    if len(users) > 0 {
        if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
            log.Fatalf("failed to insert users into mysql: %v", err)
        }
    }

    // Migrate ActiveConnections
    var conns []entities.ActiveConnection
    if err := srcDB.Find(&conns).Error; err != nil {
        log.Fatalf("failed to read active connections from sqlite: %v", err)
    }
    if len(conns) > 0 {
        if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&conns).Error; err != nil {
            log.Fatalf("failed to insert active connections into mysql: %v", err)
        }
    }

    // Migrate ControlsInformation
    var ctrs []entities.ControlsInformation
    if err := srcDB.Find(&ctrs).Error; err != nil {
        log.Fatalf("failed to read controls information from sqlite: %v", err)
    }
    if len(ctrs) > 0 {
        if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ctrs).Error; err != nil {
            log.Fatalf("failed to insert controls into mysql: %v", err)
        }
    }

    log.Printf("migration completed: users=%d conns=%d controls=%d", len(users), len(conns), len(ctrs))
}
