package router

import (
	"net/http"

	"task-api/internal/httpapi/handlers"
)

func New(h *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.Tasks)
	mux.HandleFunc("/tasks/", h.TaskByID)
	return mux
}
