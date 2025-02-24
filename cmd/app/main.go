package main

import (
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Iniciar servidor
	router := gin.Default()
	routes.SetupRoutes(router)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
