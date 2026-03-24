package main

import (
	"log"
	"net/http"

	"task-api/internal/httpapi/handlers"
	"task-api/internal/httpapi/router"
	"task-api/internal/repository"
	"task-api/internal/service"
	"task-api/internal/storage"
	userSvc "task-api/internal/userService"
)

// main собирает зависимости приложения и запускает HTTP-сервер.
func main() {
	db, err := storage.OpenDB("host=localhost user=postgres password=postgres dbname=task_api port=5432 sslmode=disable TimeZone=Europe/Warsaw")
	if err != nil {
		log.Fatal(err)
	}

	tasksRepo := repository.NewGormTaskRepository(db)
	tasksService := service.NewTaskService(tasksRepo)
	taskHandlers := handlers.NewTaskHandler(tasksService)

	userRepo := userSvc.NewGormRepository(db)
	userService := userSvc.NewService(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	mux := router.New(taskHandlers, userHandlers)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
