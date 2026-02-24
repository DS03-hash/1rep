package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

type createTaskRequest struct {
	Task   string `json:"task"`
	Status string `json:"status"`
}

type patchTaskRequest struct {
	Task   *string `json:"task"`
	Status *string `json:"status"`
}

var (
	mu     sync.Mutex
	nextID = 1
	tasks  = map[int]Task{}
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", tasksCollectionHandler)

	mux.HandleFunc("/tasks/", taskByIDHandler)

	_ = http.ListenAndServe(":8080", mux)
}

func tasksCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTask(w, r)
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
	case http.MethodGet:
		getTask(w, id)
	case http.MethodPatch:
		patchTask(w, r, id)
	case http.MethodDelete:
		deleteTask(w, id)
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
	if strings.TrimSpace(req.Status) == "" {
		req.Status = "new"
	}

	mu.Lock()
	id := nextID
	nextID++
	t := Task{ID: id, Task: req.Task, Status: req.Status}
	tasks[id] = t
	mu.Unlock()

	writeJSON(w, http.StatusCreated, t)
}

func getTask(w http.ResponseWriter, id int) {
	mu.Lock()
	t, ok := tasks[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func patchTask(w http.ResponseWriter, r *http.Request, id int) {
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	mu.Lock()
	t, ok := tasks[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if req.Task != nil {
		if strings.TrimSpace(*req.Task) == "" {
			mu.Unlock()
			http.Error(w, "task cannot be empty", http.StatusBadRequest)
			return
		}
		t.Task = *req.Task
	}
	if req.Status != nil {
		if strings.TrimSpace(*req.Status) == "" {
			mu.Unlock()
			http.Error(w, "status cannot be empty", http.StatusBadRequest)
			return
		}
		t.Status = *req.Status
	}

	tasks[id] = t
	mu.Unlock()

	writeJSON(w, http.StatusOK, t)
}

func deleteTask(w http.ResponseWriter, id int) {
	mu.Lock()
	_, ok := tasks[id]
	if ok {
		delete(tasks, id)
	}
	mu.Unlock()

	if !ok {
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
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
