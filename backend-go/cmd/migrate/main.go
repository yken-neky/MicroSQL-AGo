package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/config"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	repoport "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Django Schema Mappers
// These structs map to the actual Django tables in SQLite

// DjangoCustomUser represents Users_App_customuser
type DjangoCustomUser struct {
	ID         uint       `gorm:"column:id"`
	Password   string     `gorm:"column:password"`
	Username   string     `gorm:"column:username"`
	FirstName  string     `gorm:"column:first_name"`
	LastName   string     `gorm:"column:last_name"`
	Email      string     `gorm:"column:email"`
	IsActive   bool       `gorm:"column:is_active"`
	DateJoined time.Time  `gorm:"column:date_joined"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	LastLogin  *time.Time `gorm:"column:last_login"`
	Role       string     `gorm:"column:role"`
}

func (DjangoCustomUser) TableName() string {
	return "Users_App_customuser"
}

// DjangoConnectionLog represents Logs_App_connectionlog
type DjangoConnectionLog struct {
	ID        uint      `gorm:"column:id"`
	Driver    string    `gorm:"column:driver"`
	Server    string    `gorm:"column:server"`
	DBUser    string    `gorm:"column:db_user"`
	Timestamp time.Time `gorm:"column:timestamp"`
	Status    string    `gorm:"column:status"`
	UserID    uint      `gorm:"column:user_id"`
}

func (DjangoConnectionLog) TableName() string {
	return "Logs_App_connectionlog"
}

// DjangoActiveConnection represents Connecting_App_activeconnection
type DjangoActiveConnection struct {
	ID            uint      `gorm:"column:id"`
	Driver        string    `gorm:"column:driver"`
	Server        string    `gorm:"column:server"`
	DBUser        string    `gorm:"column:db_user"`
	Password      string    `gorm:"column:password"`
	IsConnected   bool      `gorm:"column:is_connected"`
	LastConnected time.Time `gorm:"column:last_connected"`
	UserID        uint      `gorm:"column:user_id"`
}

func (DjangoActiveConnection) TableName() string {
	return "Connecting_App_activeconnection"
}

// DjangoControlsInformation represents InsideDB_App_controls_information
type DjangoControlsInformation struct {
	ID          uint   `gorm:"column:id"`
	Idx         int    `gorm:"column:idx"`
	Name        string `gorm:"column:name"`
	Chapter     string `gorm:"column:chapter"`
	Description string `gorm:"column:description"`
	Impact      string `gorm:"column:impact"`
	GoodConfig  string `gorm:"column:good_config"`
	BadConfig   string `gorm:"column:bad_config"`
	Ref         string `gorm:"column:ref"`
}

func (DjangoControlsInformation) TableName() string {
	return "InsideDB_App_controls_information"
}

// Migration entrypoint: reads Django tables from SQLite and maps to Go entities in MySQL
func main() {
	cfg := config.LoadConfig()

	srcPath := flag.String("src", cfg.DBPath, "path to source sqlite database file (overrides DB_PATH)")
	dstDSN := flag.String("dst-dsn", "", "MySQL DSN (user:pass@tcp(host:port)/dbname?params). If empty, uses MYSQL_* env vars from Config")
	flag.Parse()

	// Open source (sqlite with Django schema)
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

	// Ensure destination has tables for all domain entities
	log.Println("Creating tables in destination MySQL...")
	if err := dstDB.AutoMigrate(
		&entities.User{},
		&entities.ActiveConnection{},
		&entities.ControlsInformation{},
		&repoport.ControlsScript{},
		&entities.Role{},
		&entities.Permission{},
		&entities.UserRole{},
		&entities.ConnectionLog{},
		&entities.AdminActionLog{},
		// queries and related result/stats were removed from the Go model (user-facing SQL execution removed)
	); err != nil {
		log.Fatalf("failed to auto-migrate destination DB: %v", err)
	}

	// Seed default roles and permissions in destination DB as part of migration
	if err := repositories.SeedDefaultRolesAndPermissions(dstDB); err != nil {
		log.Fatalf("failed to seed roles/permissions on destination DB: %v", err)
	}

	// Migration counters
	var (
		usersCopied    int
		connsCopied    int
		connLogsCopied int
		controlsCopied int
		scriptsCopied  int
	)

	// 1. Migrate Users from Django CustomUser to Go User
	log.Println("Migrating users...")
	var djangoUsers []DjangoCustomUser
	if err := srcDB.Find(&djangoUsers).Error; err != nil {
		log.Printf("warning: failed to read users from sqlite: %v (continuing...)\n", err)
	} else if len(djangoUsers) > 0 {
		goUsers := make([]entities.User, len(djangoUsers))
		for i, du := range djangoUsers {
			goUsers[i] = entities.User{
				ID:        du.ID,
				Username:  du.Username,
				Email:     du.Email,
				Password:  du.Password,
				FirstName: du.FirstName,
				LastName:  du.LastName,
				Role:      du.Role,
				CreatedAt: du.CreatedAt,
				IsActive:  du.IsActive,
			}
			if du.LastLogin != nil {
				// entities.User.LastLogin is nullable (*time.Time) so assign pointer directly
				goUsers[i].LastLogin = du.LastLogin
			}
		}
		if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&goUsers).Error; err != nil {
			log.Fatalf("failed to insert users into mysql: %v", err)
		}
		usersCopied = len(goUsers)
		log.Printf("✓ Migrated %d users\n", usersCopied)
	}

	// 2. Migrate ActiveConnections
	log.Println("Migrating active connections...")
	var djangoConns []DjangoActiveConnection
	if err := srcDB.Find(&djangoConns).Error; err != nil {
		log.Printf("warning: failed to read active connections from sqlite: %v (continuing...)\n", err)
	} else if len(djangoConns) > 0 {
		goConns := make([]entities.ActiveConnection, len(djangoConns))
		for i, dc := range djangoConns {
			// Filter out zero/null dates (MySQL rejects '0000-00-00')
			lastConnected := dc.LastConnected
			if lastConnected.IsZero() {
				lastConnected = time.Now()
			}
			goConns[i] = entities.ActiveConnection{
				ID:     dc.ID,
				UserID: dc.UserID,
				// older schema stores driver string; populate Manager from Driver for compatibility
				Driver:        dc.Driver,
				Manager:       dc.Driver,
				Server:        dc.Server,
				DBUser:        dc.DBUser,
				Password:      dc.Password,
				IsConnected:   dc.IsConnected,
				LastConnected: lastConnected,
			}
		}
		if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&goConns).Error; err != nil {
			log.Fatalf("failed to insert active connections into mysql: %v", err)
		}
		connsCopied = len(goConns)
		log.Printf("✓ Migrated %d active connections\n", connsCopied)
	}

	// 3. Migrate ConnectionLogs
	log.Println("Migrating connection logs...")
	var djangoLogs []DjangoConnectionLog
	if err := srcDB.Find(&djangoLogs).Error; err != nil {
		log.Printf("warning: failed to read connection logs from sqlite: %v (continuing...)\n", err)
	} else if len(djangoLogs) > 0 {
		goLogs := make([]entities.ConnectionLog, len(djangoLogs))
		for i, dl := range djangoLogs {
			goLogs[i] = entities.ConnectionLog{
				ID:        dl.ID,
				UserID:    dl.UserID,
				Driver:    dl.Driver,
				Server:    dl.Server,
				DBUser:    dl.DBUser,
				Timestamp: dl.Timestamp,
				Status:    dl.Status,
			}
		}
		if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&goLogs).Error; err != nil {
			log.Fatalf("failed to insert connection logs into mysql: %v", err)
		}
		connLogsCopied = len(goLogs)
		log.Printf("✓ Migrated %d connection logs\n", connLogsCopied)
	}

	// 4. Migrate ControlsInformation
	log.Println("Migrating controls information...")
	var djangoControls []DjangoControlsInformation
	if err := srcDB.Find(&djangoControls).Error; err != nil {
		log.Printf("warning: failed to read controls information from sqlite: %v (continuing...)\n", err)
	} else if len(djangoControls) > 0 {
		goControls := make([]entities.ControlsInformation, len(djangoControls))
		for i, dc := range djangoControls {
			goControls[i] = entities.ControlsInformation{
				ID:          dc.ID,
				Idx:         dc.Idx,
				Name:        dc.Name,
				Chapter:     dc.Chapter,
				Description: dc.Description,
			}
		}
		if err := dstDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&goControls).Error; err != nil {
			log.Fatalf("failed to insert controls into mysql: %v", err)
		}
		controlsCopied = len(goControls)
		log.Printf("✓ Migrated %d controls information\n", controlsCopied)
	}

	// 5. Migrate Control Scripts (handle renamed source table/column)
	log.Println("Migrating control scripts...")

	// try common table name variants in source sqlite
	tableCandidates := []string{"controls_script", "InsideDB_App_controls_scripts"}
	var srcTable string
	for _, t := range tableCandidates {
		if srcDB.Migrator().HasTable(t) {
			srcTable = t
			break
		}
	}

	if srcTable == "" {
		log.Printf("warning: no control scripts table found in source (tried %v) - skipping scripts migration\n", tableCandidates)
	} else {
		type DjangoControlsScript struct {
			ID               uint   `gorm:"column:id"`
			ControlType      string `gorm:"column:control_type"`
			QuerySQL         string `gorm:"column:query_sql"`
			ControlScriptRef uint   `gorm:"column:control_script_ref"`
		}

		var djangoScripts []DjangoControlsScript

		// try possible fk column names and stop on first that yields rows
		fkCandidates := []string{"control_script_id", "control_script_id_id"}
		for _, fk := range fkCandidates {
			selectQuery := fmt.Sprintf("SELECT id, control_type, query_sql, %s AS control_script_ref FROM %s", fk, srcTable)
			if err := srcDB.Raw(selectQuery).Scan(&djangoScripts).Error; err != nil {
				log.Printf("debug: query failed for fk='%s' table='%s': %v\n", fk, srcTable, err)
				continue
			}
			if len(djangoScripts) > 0 {
				// map and insert into controls_scripts table with normalized column name
				goScripts := make([]repoport.ControlsScript, len(djangoScripts))
				for i, ds := range djangoScripts {
					goScripts[i] = repoport.ControlsScript{
						ID:               ds.ID,
						ControlType:      ds.ControlType,
						QuerySQL:         ds.QuerySQL,
						ControlScriptRef: ds.ControlScriptRef,
					}
				}
				// Write to controls_scripts table (GORM will use TableName() from ControlsScript)
				if err := dstDB.Table("controls_scripts").Clauses(clause.OnConflict{DoNothing: true}).Create(&goScripts).Error; err != nil {
					log.Fatalf("failed to insert control scripts into mysql: %v", err)
				}
				scriptsCopied = len(goScripts)
				log.Printf("✓ Migrated %d control scripts (using fk '%s' from table '%s')\n", scriptsCopied, fk, srcTable)
				break
			}
		}
	}
	// Summary
	log.Printf("\n✓ Migration completed successfully!\n")
	log.Printf("  Users:               %d\n", usersCopied)
	log.Printf("  Active Connections:  %d\n", connsCopied)
	log.Printf("  Connection Logs:     %d\n", connLogsCopied)
	log.Printf("  Controls:            %d\n", controlsCopied)
	log.Printf("  Control Scripts:     %d\n", scriptsCopied)
}
