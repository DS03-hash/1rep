DB_DSN := postgres://postgres:postgres@localhost:5432/task_api?sslmode=disable
MIGRATE := migrate -path ./migrations -database "$(DB_DSN)"
GOLANGCI_LINT_VERSION := v1.64.8
GOLANGCI_LINT := $(shell go env GOPATH)\bin\golangci-lint.exe

.PHONY: migrate-new migrate migrate-down migrate-down-one run gen lint-install lint test check

# Создать новую миграцию: make migrate-new NAME=add_tasks_table
migrate-new:
	migrate create -ext sql -dir ./migrations $(NAME)

# Применить все новые миграции.
migrate:
	$(MIGRATE) up

# Откатить все применённые миграции.
migrate-down:
	$(MIGRATE) down

# Откатить одну последнюю миграцию.
migrate-down-one:
	$(MIGRATE) down 1

# Запустить API-сервис.
run:
	go run cmd/api/main.go

# Сгенерировать типы и интерфейсы из OpenAPI.
gen:
	if not exist internal\httpapi\gen mkdir internal\httpapi\gen
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 -config oapi-codegen.yaml openapi.yaml

# Установить golangci-lint локально.
lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# Запустить линтеры.
lint:
	"$(GOLANGCI_LINT)" run ./...

# Запустить go-тесты.
test:
	go test ./...

# Запустить все локальные проверки перед push.
check: test lint
