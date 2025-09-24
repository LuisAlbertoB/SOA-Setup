package models

import "gorm.io/gorm"

// User representa la tabla 'users' en la base de datos.
type User struct {
	gorm.Model        // Incluye ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"size:255;not null"`
	Email      string `gorm:"size:255;not null;uniqueIndex"`
	Password   string `gorm:"size:255;not null"`
	Phone      string `gorm:"size:50"`
	Region     string `gorm:"size:100"`
	// Relación: Un usuario tiene muchas Apps
	Apps []App `gorm:"foreignKey:DeveloperID"`
}
