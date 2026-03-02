package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-api/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type createTaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

type patchTaskRequest struct {
	Task   *string `json:"task"`
	IsDone *bool   `json:"is_done"`
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.listTasks(w)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	uid := uint(id)

	switch r.Method {
	case http.MethodPatch:
		h.patchTask(w, r, uid)
	case http.MethodPut:
		h.patchTask(w, r, uid) // минимально: PUT как PATCH
	case http.MethodDelete:
		h.deleteTask(w, uid)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	t, err := h.svc.Create(req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, "task is required", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *TaskHandler) listTasks(w http.ResponseWriter) {
	tasks, err := h.svc.List()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) patchTask(w http.ResponseWriter, r *http.Request, id uint) {
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	t, err := h.svc.Patch(id, req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, id uint) {
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
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
