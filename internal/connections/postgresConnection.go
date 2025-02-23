package connections

import (
	"log"
	"main/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	dsn := `host=localhost 
			user=postgres 
			password=POSTGRE*SQL 
			dbname=MicroSQL_AGo 
			port=5432 
			sslmode=disable`

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error de conexión:", err)
	}

	// Configurar pool de conexiones
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Migración
	if err := models.AutoMigrate(DB); err != nil {
		log.Fatal("Error en migración:", err)
	}

	log.Println("¡BD conectada y migrada!")
}
