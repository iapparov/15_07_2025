package web

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, userHandler *UserHandler) {
	mux.HandleFunc("/tasks", userHandler.CreateTask)
	mux.HandleFunc("/tasks/", userHandler.Tasks)
}

