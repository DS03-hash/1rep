package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"task-api/internal/httpapi/gen"
	"task-api/internal/userService"
)

type UserHandlers struct {
	svc *userService.Service
}

// NewUserHandlers creates HTTP handler for users.
func NewUserHandlers(svc *userService.Service) *UserHandlers {
	return &UserHandlers{svc: svc}
}

// Users routes requests to user collection endpoints.
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

// UserByID routes requests for single user by id.
func (h *UserHandlers) UserByID(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandlers) listUsers(w http.ResponseWriter) {
	users, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, toListUsersResponse(users))
}

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

func toCreateUserResponse(u userService.User) gen.CreateUserResponse {
	return gen.CreateUserResponse{
		Id:        int64(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toPatchUserResponse(u userService.User) gen.PatchUserResponse {
	return gen.PatchUserResponse{
		Id:        int64(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

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
