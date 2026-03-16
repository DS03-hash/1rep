package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-api/internal/domain"
	"task-api/internal/httpapi/gen"
	"task-api/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

// NewTaskHandler создает HTTP-обработчик задач.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// Tasks маршрутизирует запросы к коллекции задач по HTTP-методу.
func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.listTasks(w)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// TaskByID маршрутизирует запросы к одной задаче по ее идентификатору.
func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	uid := uint(id)

	switch r.Method {
	case http.MethodPatch:
		h.patchTask(w, r, uid)
	case http.MethodDelete:
		h.deleteTask(w, uid)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// createTask декодирует тело запроса и создает задачу.
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req gen.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	t, err := h.svc.Create(req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "task is required")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, toCreateTaskResponse(*t))
}

// listTasks возвращает список всех задач.
func (h *TaskHandler) listTasks(w http.ResponseWriter) {
	tasks, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, toListTasksResponse(tasks))
}

// patchTask частично обновляет задачу по id.
func (h *TaskHandler) patchTask(w http.ResponseWriter, r *http.Request, id uint) {
	var req gen.PatchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	t, err := h.svc.Patch(id, req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPatchTaskResponse(*t))
}

// deleteTask удаляет задачу по id.
func (h *TaskHandler) deleteTask(w http.ResponseWriter, id uint) {
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseIDFromPath извлекает числовой id из URL-пути.
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

// toCreateTaskResponse преобразует доменную задачу в ответ create.
func toCreateTaskResponse(t domain.Task) gen.CreateTaskResponse {
	return gen.CreateTaskResponse{
		Id:     int64(t.ID),
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

// toPatchTaskResponse преобразует доменную задачу в ответ patch.
func toPatchTaskResponse(t domain.Task) gen.PatchTaskResponse {
	return gen.PatchTaskResponse{
		Id:     int64(t.ID),
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

// toListTasksResponse преобразует список доменных задач в API-ответ.
func toListTasksResponse(tasks []domain.Task) gen.ListTasksResponse {
	out := make(gen.ListTasksResponse, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, gen.Task{
			Id:     int64(t.ID),
			Task:   t.Task,
			IsDone: t.IsDone,
		})
	}
	return out
}

// writeJSON отправляет JSON-ответ с указанным статусом.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError отправляет стандартную ошибку API в JSON-формате.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, gen.ErrorResponse{Message: message})
}
