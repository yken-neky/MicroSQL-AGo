package connections

import (
	"log"
	"main/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabase initializes and returns a connection to the PostgreSQL database using GORM.
// It sets up the connection pool and performs database migration.
//
// Returns:
//   *gorm.DB: A pointer to the GORM database instance.
//
// The function performs the following steps:
// 1. Constructs the DSN (Data Source Name) for the PostgreSQL connection.
// 2. Attempts to open a connection to the database using GORM.
// 3. Logs a fatal error and exits if the connection fails.
// 4. Logs a success message if the connection is established.
// 5. Configures the connection pool with a maximum of 10 idle connections and 100 open connections.
// 6. Performs database migration using the AutoMigrate function from the models package.
// 7. Logs a fatal error and exits if the migration fails.
// 8. Logs a success message if the migration is successful.
// 9. Returns the GORM database instance.
//
// Note: The function logs errors and exits the application if any step fails.
func SetupDatabase() *gorm.DB {
	dsn := `host=localhost 
			user=postgres 
			password=POSTGRE*SQL 
			dbname=MicroSQL_AGo 
			port=5432 
			sslmode=disable`
	 
	if DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatal("Error de conexión", err)
	} else {
		log.Println("¡BD conectada!", err)

		// Configurar pool de conexiones
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		// Migración
		if err := models.AutoMigrate(DB); err != nil {
			log.Fatal("Error en migración")
		}

		log.Println("¡BD conectada y migrada!")
		return DB
	}
	return nil
}


	