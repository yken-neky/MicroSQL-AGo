package main

import (
	"main/local/connections"
	"main/local/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar BD
	connections.SetupDatabase()

	// Iniciar servidor
	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
