package main

import (
	"main/internal/connections"
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar BD
	connections.SetupDatabase()

	// Iniciar servidor
	router := gin.Default()
	routes.SetupRoutes(router)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
