package router

import (
	"net/http"
	"rob-api-go/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/s3/get-app-files/{appId}", handlers.GetAppFilesHandler)
	r.Post("/s3/upload-app-files", handlers.UploadAppFilesHandler)

	return r
}
