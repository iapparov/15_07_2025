package web

import (
	"net/http"
	"zip_downloader/internal/config"
	"zip_downloader/internal/datasource"
	"github.com/google/uuid"
	"encoding/json"
)

type UserHandler struct{
	Repo map[uuid.UUID]datasource.Storage
	Config *config.Config
}

func NewUserHandler(conf *config.Config) *UserHandler{
	return &UserHandler{
		Repo: make(map[uuid.UUID]datasource.Storage),
		Config: conf,
	}
}

func (h *UserHandler) CreateTask(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	newuuid := uuid.New()
	data := datasource.Storage{
		Id: newuuid,
		Files: make([]string, 0, h.Config.Max_objects),
		Archive: "",
		Status: datasource.StatusWait,
	}
	h.Repo[newuuid]=data

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newuuid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task created\n"))

}

func (h *UserHandler) Task(w http.ResponseWriter, r *http.Request){

}

