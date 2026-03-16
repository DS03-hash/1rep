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
// New - это функция, которая создает новый HTTP-маршрутизатор (ServeMux) и регистрирует обработчики для путей "/tasks" и "/tasks/{id}".
// Она принимает ссылку на TaskHandler и возвращает настроенный ServeMux, который будет использоваться для обработки входящих HTTP-запросов. 
// Путь "/tasks" обрабатывается методом Tasks, а путь "/tasks/{id}" обрабатывается методом TaskByID.
