package models

import (
	"time"

	"gorm.io/gorm"
)

// App representa la tabla 'apps'
type App struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Version     string `gorm:"size:50;not null"`
	ReleaseDate time.Time
	DeveloperID uint // Clave foránea para la relación con User
}
