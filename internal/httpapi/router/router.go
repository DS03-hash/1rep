package router

import (
	"net/http"

	"task-api/internal/httpapi/handlers"
)

// New configures HTTP routes for tasks and users API.
func New(taskHandlers *handlers.TaskHandler, userHandlers *handlers.UserHandlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", taskHandlers.Tasks)
	mux.HandleFunc("/tasks/", taskHandlers.TaskByID)
	mux.HandleFunc("/users", userHandlers.Users)
	mux.HandleFunc("/users/", userHandlers.UserByID)
	return mux
}
