package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"rob-api-go/internal/database"
	"rob-api-go/internal/models"
	"rob-api-go/internal/services"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// Helper para enviar respuestas JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Helper para enviar errores JSON
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// GetAppFilesHandler obtiene las URLs firmadas para los archivos de una app.
func GetAppFilesHandler(w http.ResponseWriter, r *http.Request) {
	appIDStr := chi.URLParam(r, "appId")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de aplicación inválido")
		return
	}

	var appFile models.AppFile
	if err := database.DB.Preload("Screenshots").First(&appFile, "app_id = ?", appID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Archivos de la aplicación no encontrados")
		return
	}

	s3Svc := services.NewS3Service()

	iconURL, _ := s3Svc.GetSignedURL(appFile.IconKey)
	appFileURL, _ := s3Svc.GetSignedURL(appFile.AppFileKey)

	screenshotURLs := make([]string, len(appFile.Screenshots))
	for i, ss := range appFile.Screenshots {
		url, _ := s3Svc.GetSignedURL(ss.S3Key)
		screenshotURLs[i] = url
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"iconUrl":     iconURL,
		"appFileUrl":  appFileURL,
		"screenshots": screenshotURLs,
	})
}

// UploadAppFilesHandler maneja la subida de archivos de una aplicación.
func UploadAppFilesHandler(w http.ResponseWriter, r *http.Request) {
	// El tamaño máximo del formulario es de 100MB
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "No se pudo parsear el formulario multipart")
		return
	}

	appIDStr := r.FormValue("appId")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "appId es requerido y debe ser un número")
		return
	}

	// Aquí iría la lógica de autorización para verificar que el usuario es el dueño de la app
	// Por simplicidad, la omitimos en este ejemplo.

	s3Svc := services.NewS3Service()
	var appFile models.AppFile
	database.DB.FirstOrCreate(&appFile, models.AppFile{AppID: uint(appID)})

	// Manejar subida de ícono
	icon, iconHeader, err := r.FormFile("icon")
	if err == nil {
		defer icon.Close()
		if appFile.IconKey != "" {
			s3Svc.DeleteFile(appFile.IconKey)
		}
		key := fmt.Sprintf("apps/%d/icon%s", appID, filepath.Ext(iconHeader.Filename))
		_, err := s3Svc.UploadFile(iconHeader, key)
		if err != nil {
			log.Printf("Error al subir ícono: %v", err)
			respondError(w, http.StatusInternalServerError, "Error al subir el ícono")
			return
		}
		appFile.IconKey = key
	}

	// Manejar subida de archivo de la app (APK)
	app, appHeader, err := r.FormFile("appFile")
	if err == nil {
		defer app.Close()
		if appFile.AppFileKey != "" {
			s3Svc.DeleteFile(appFile.AppFileKey)
		}
		key := fmt.Sprintf("apps/%d/app%s", appID, filepath.Ext(appHeader.Filename))
		_, err := s3Svc.UploadFile(appHeader, key)
		if err != nil {
			log.Printf("Error al subir app: %v", err)
			respondError(w, http.StatusInternalServerError, "Error al subir el archivo de la app")
			return
		}
		appFile.AppFileKey = key
	}

	// Manejar subida de capturas de pantalla
	screenshots := r.MultipartForm.File["screenshots"]
	if len(screenshots) > 0 {
		// Lógica para borrar capturas viejas (omitida por simplicidad)
		// ...

		for _, ssHeader := range screenshots {
			key := fmt.Sprintf("apps/%d/screenshots/%d%s", appID, time.Now().UnixNano(), filepath.Ext(ssHeader.Filename))
			_, err := s3Svc.UploadFile(ssHeader, key)
			if err != nil {
				log.Printf("Error al subir screenshot: %v", err)
				continue // Continuar con la siguiente
			}
			newScreenshot := models.Screenshot{
				AppFileID: appFile.ID,
				S3Key:     key,
			}
			database.DB.Create(&newScreenshot)
		}
	}

	if err := database.DB.Save(&appFile).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Error al guardar los datos de los archivos")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Archivos subidos correctamente"})
}

// NOTA: El handler para `DeleteAppFiles` seguiría una lógica similar:
// 1. Verificar autorización.
// 2. Obtener todas las keys de la BD.
// 3. Llamar a `s3Svc.DeleteFile` para cada key.
// 4. Eliminar los registros de la BD.
