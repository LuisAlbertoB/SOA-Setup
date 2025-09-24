package database

import (
	"fmt"
	"log"
	"rob-api-go/internal/config"
	"rob-api-go/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect se conecta a la base de datos y ejecuta las migraciones.
func Connect() {
	dsn := config.Config.MariaDBDSN
	if dsn == "" {
		log.Fatal("Error: La variable de entorno MARIADB_DSN no está definida.")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	fmt.Println("Conectado a la base de datos.")

	// AutoMigrate crea/actualiza las tablas según los modelos
	err = DB.AutoMigrate(&models.User{}, &models.App{}, &models.AppFile{}, &models.Screenshot{})
	if err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}
	fmt.Println("Migraciones completadas.")
}
