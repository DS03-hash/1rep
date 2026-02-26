package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type createTaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

type patchTaskRequest struct {
	Task   *string `json:"task"`
	IsDone *bool   `json:"is_done"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&Task{}); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", tasksCollectionHandler)
	mux.HandleFunc("/tasks/", taskByIDHandler)

	_ = http.ListenAndServe(":8080", mux)
}

func tasksCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTask(w, r)
	case http.MethodGet:
		listTasks(w)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		patchTask(w, r, uint(id))
	case http.MethodPut:

		patchTask(w, r, uint(id))
	case http.MethodDelete:
		deleteTask(w, uint(id))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Task) == "" {
		http.Error(w, "task is required", http.StatusBadRequest)
		return
	}

	t := Task{
		Task:   req.Task,
		IsDone: req.IsDone,
	}

	if err := db.Create(&t).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func listTasks(w http.ResponseWriter) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func patchTask(w http.ResponseWriter, r *http.Request, id uint) {
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	var t Task
	if err := db.First(&t, id).Error; err != nil {

		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if req.Task != nil {
		if strings.TrimSpace(*req.Task) == "" {
			http.Error(w, "task cannot be empty", http.StatusBadRequest)
			return
		}
		t.Task = *req.Task
	}
	if req.IsDone != nil {
		t.IsDone = *req.IsDone
	}

	if err := db.Save(&t).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func deleteTask(w http.ResponseWriter, id uint) {

	res := db.Delete(&Task{}, id)
	if res.Error != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if res.RowsAffected == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDFromPath(path, prefix string) (int, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, errors.New("bad path")
	}
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, errors.New("empty id")
	}
	return strconv.Atoi(raw)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
