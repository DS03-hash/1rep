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
	// подключение к БД с помощью функции OpenDB из пакета storage. В случае ошибки подключения,
	// программа завершится с сообщением об ошибке.

	repo := repository.NewGormTaskRepository(db)
	// создание нового репозитория задач с помощью функции NewGormTaskRepository из пакета repository,
	// которая принимает подключение к базе данных.
	svc := service.NewTaskService(repo)
	// создание нового сервиса задач с помощью функции NewTaskService из пакета service,
	// которая принимает репозиторий задач. Сервис будет использоваться для обработки бизнес-логики, связанной с задачами.
	h := handlers.NewTaskHandler(svc)
	// создание нового обработчика задач с помощью функции NewTaskHandler из пакета handlers,
	// которая принимает сервис задач. Обработчик будет использоваться для обработки HTTP-запросов, связанных с задачами.

	mux := router.New(h)
	// создание нового маршрутизатора с помощью функции New из пакета router,
	// которая принимает обработчик задач. Маршрутизатор будет использоваться для маршрутизации HTTP-запросов к
	//  соответствующим обработчикам.

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
	// запуск HTTP-сервера на порту 8080 с использованием маршрутизатора mux.
	// Если сервер не может запуститься, программа завершится с сообщением об ошибке.
}
