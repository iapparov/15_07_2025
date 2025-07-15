package web

import (
	"net/http"
	"zip_downloader/internal/config"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, userHandler *UserHandler, config config.Config) {
	r.Post("/tasks", userHandler.CreateTask)
	r.Post("/tasks/{id}/add", userHandler.AddToTask)
	r.Get("/tasks/{id}/status", userHandler.TaskStatus)
	r.Handle("/archives/*", http.StripPrefix("/archives/", http.FileServer(http.Dir(config.Archive_path))))
}