package main

import (
	"log"
	"net/http"

	"task-api/internal/httpapi/handlers"
	"task-api/internal/httpapi/router"
	"task-api/internal/repository"
	"task-api/internal/service"
	"task-api/internal/storage"
)

func main() {
	db, err := storage.OpenDB("host=localhost user=postgres password=postgres dbname=task_api port=5432 sslmode=disable TimeZone=Europe/Warsaw")
	if err != nil {
		log.Fatal(err)
	}
	// Подлкючение к БД

	repo := repository.NewGormTaskRepository(db)
	// Инициализация репозитория

	svc := service.NewTaskService(repo)
	// Инициализация сервиса

	h := handlers.NewTaskHandler(svc)
	// Инициализация HTTP-обработчика

	mux := router.New(h)
	// Инициализация маршрутизатора

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
	// Запуск HTTP-сервера
}
