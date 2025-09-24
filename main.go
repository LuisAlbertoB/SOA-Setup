package main

import (
	"fmt"
	"log"
	"net/http"
	"rob-api-go/internal/config"
	"rob-api-go/internal/database"
	"rob-api-go/internal/router"
)

func main() {
	// Cargar la configuración centralizada
	config.LoadConfig()

	// Conectar a la BD y migrar
	database.Connect()

	// Configurar el router
	r := router.SetupRouter()

	log.Printf("Servidor iniciado en http://localhost:%s\n", config.Config.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Config.ServerPort), r))
}
