package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-api/internal/httpapi/gen"
	"task-api/internal/userService"
)

type UserHandlers struct {
	svc *userService.Service
}

// NewUserHandlers создает HTTP-обработчик пользователей.
func NewUserHandlers(svc *userService.Service) *UserHandlers {
	return &UserHandlers{svc: svc}
}

// Users маршрутизирует запросы к коллекции пользователей по HTTP-методу.
func (h *UserHandlers) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodGet:
		h.listUsers(w)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// UserByID маршрутизирует запросы к одному пользователю по его идентификатору.
func (h *UserHandlers) UserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/tasks") {
		h.listTasksByUserID(w, r)
		return
	}

	id, err := parseIDFromPath(r.URL.Path, "/users/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	uid := uint(id)

	switch r.Method {
	case http.MethodPatch:
		h.patchUser(w, r, uid)
	case http.MethodDelete:
		h.deleteUser(w, uid)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// listTasksByUserID возвращает задачи конкретного пользователя.
func (h *UserHandlers) listTasksByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserIDForTasksPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tasks, err := h.svc.GetTasksForUser(uint(id))
	if err != nil {
		if errors.Is(err, userService.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toListTasksResponse(tasks))
}

// createUser декодирует тело запроса и создает пользователя.
func (h *UserHandlers) createUser(w http.ResponseWriter, r *http.Request) {
	var req gen.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	u, err := h.svc.Create(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, userService.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, toCreateUserResponse(*u))
}

// listUsers возвращает список всех пользователей.
func (h *UserHandlers) listUsers(w http.ResponseWriter) {
	users, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, toListUsersResponse(users))
}

// patchUser частично обновляет пользователя по id.
func (h *UserHandlers) patchUser(w http.ResponseWriter, r *http.Request, id uint) {
	var req gen.PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	u, err := h.svc.Patch(id, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, userService.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, userService.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPatchUserResponse(*u))
}

// deleteUser удаляет пользователя по id.
func (h *UserHandlers) deleteUser(w http.ResponseWriter, id uint) {
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, userService.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// toCreateUserResponse преобразует модель пользователя в ответ create.
func toCreateUserResponse(u userService.User) gen.CreateUserResponse {
	return gen.CreateUserResponse{
		Id:        int64(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toPatchUserResponse преобразует модель пользователя в ответ patch.
func toPatchUserResponse(u userService.User) gen.PatchUserResponse {
	return gen.PatchUserResponse{
		Id:        int64(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toListUsersResponse преобразует список пользователей в API-ответ.
func toListUsersResponse(users []userService.User) gen.ListUsersResponse {
	out := make(gen.ListUsersResponse, 0, len(users))
	for _, u := range users {
		out = append(out, gen.User{
			Id:        int64(u.ID),
			Email:     u.Email,
			Password:  u.Password,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return out
}

func parseUserIDForTasksPath(path string) (int, error) {
	const (
		prefix = "/users/"
		suffix = "/tasks"
	)

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return 0, errors.New("bad path")
	}

	raw := strings.TrimPrefix(path, prefix)
	raw = strings.TrimSuffix(raw, suffix)
	raw = strings.Trim(raw, "/")
	if raw == "" || strings.Contains(raw, "/") {
		return 0, errors.New("invalid id")
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}

	return id, nil
}
