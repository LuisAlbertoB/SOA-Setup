package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig almacena toda la configuración de la aplicación.
type AppConfig struct {
	ServerPort      string
	MariaDBDSN      string
	AWSRegion       string
	AWSS3BucketName string
}

// Config es una instancia global de la configuración de la app.
var Config AppConfig

// LoadConfig carga la configuración desde un archivo .env o desde el entorno.
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, usando variables de entorno del sistema.")
	}

	Config = AppConfig{
		ServerPort:      os.Getenv("SERVER_PORT"),
		MariaDBDSN:      os.Getenv("MARIADB_DSN"),
		AWSRegion:       os.Getenv("AWS_REGION"),
		AWSS3BucketName: os.Getenv("AWS_S3_BUCKET_NAME"),
	}

	// Aquí se podrían añadir validaciones, por ejemplo, para asegurar que las variables no estén vacías.
}
