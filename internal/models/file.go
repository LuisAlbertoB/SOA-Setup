package models

import "gorm.io/gorm"

// AppFile representa los archivos asociados a una App.
// Se establece una relación uno a uno con la tabla App.
type AppFile struct {
	gorm.Model
	AppID      uint   `gorm:"uniqueIndex"` // Relación uno a uno con App
	IconKey    string `gorm:"size:512"`    // Clave del objeto en S3
	AppFileKey string `gorm:"size:512"`    // Clave del APK en S3
	// Relación uno a muchos con capturas de pantalla
	Screenshots []Screenshot `gorm:"foreignKey:AppFileID"`
}

// Screenshot representa una captura de pantalla para una App
type Screenshot struct {
	gorm.Model
	AppFileID uint   `gorm:"not null"`
	S3Key     string `gorm:"size:512;not null"`
}
