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

	repo := repository.NewGormTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handlers.NewTaskHandler(svc)

	mux := router.New(h)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
