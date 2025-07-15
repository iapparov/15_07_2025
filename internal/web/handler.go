package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"zip_downloader/internal/app"
	"zip_downloader/internal/config"
	"zip_downloader/internal/datasource"
	"sync"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct{
	Repo map[uuid.UUID]datasource.Storage
	Config *config.Config
	mu sync.Mutex
}

func NewUserHandler(conf *config.Config) *UserHandler{
	return &UserHandler{
		Repo: make(map[uuid.UUID]datasource.Storage),
		Config: conf,
	}
}

func (h *UserHandler) CreateTask(w http.ResponseWriter, r *http.Request){
	
	if h.activeTasksCount() >= h.Config.Max_tasks {
		http.Error(w, "server busy: maximum number of active tasks reached", http.StatusTooManyRequests)
		return
	}
	data := datasource.Storage{
		Id: uuid.New(),
		Files: make([]string, 0, h.Config.Max_objects),
		Archive: "",
		Status: datasource.StatusWait,
	}
	h.Repo[data.Id]=data
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data.Id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task created\n"))

}

func (h *UserHandler) AddToTask(w http.ResponseWriter, r *http.Request){
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	storage, exists := h.Repo[taskID]
	if !exists {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	existing := len(storage.Files)

	var added, skipped []string

	for _, link := range body {
		if existing >= h.Config.Max_objects {
			skipped = append(skipped, link+" (task full)")
			continue
		}

		ext := strings.ToLower(filepath.Ext(link))
		if !h.Config.AllowedTypesMap[ext] {
			skipped = append(skipped, link+" (invalid file type)")
			continue
		}

		storage.Files = append(storage.Files, link)
		added = append(added, link)
		existing++
	}

	if existing == 3{
		storage.Status = datasource.StatusReady
	}

	// Обновляем задачу
	h.Repo[taskID] = storage

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"added":   added,
		"skipped": skipped,
	})
}

func (h *UserHandler) TaskStatus(w http.ResponseWriter, r *http.Request){
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	_, exists := h.Repo[taskID]
	if !exists {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("task status: "))
	json.NewEncoder(w).Encode(h.Repo[taskID].Status)
	if h.Repo[taskID].Status == datasource.StatusDone {
		fmt.Fprintf(w, "archive address: http://localhost:%d/archives/%s\n", h.Config.Http_port, h.Repo[taskID].Archive)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) StartWorker() {
    go func() {
        for {
			h.mu.Lock()
            for id, task := range h.Repo {
                if task.Status == datasource.StatusReady && task.Archive == "" {
                    err := zip_service.DownloadAndArchive(task.Files, h.Config.Archive_path, task.Id)
                    if err != nil {
						fmt.Println(err)
                        task.Status = datasource.StatusError
                    } else {
                        task.Status = datasource.StatusDone
						task.Archive = fmt.Sprintf("%s.zip", task.Id.String())
                    }
                }
				 h.Repo[id] = task
            }
			h.mu.Unlock()
            time.Sleep(2 * time.Second)
        }
    }()
}

func (h *UserHandler) activeTasksCount() int {
	count := 0
	for _, task := range h.Repo {
		if task.Status == datasource.StatusWait {
			count++
		}
	}
	return count
}