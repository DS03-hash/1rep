DB_DSN := postgres://postgres:postgres@localhost:5432/task_api?sslmode=disable
MIGRATE := migrate -path ./migrations -database "$(DB_DSN)"
GOLANGCI_LINT_VERSION := v1.64.8
GOLANGCI_LINT := $(shell go env GOPATH)\bin\golangci-lint.exe

migrate-new:
	migrate create -ext sql -dir ./migrations $(NAME)

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 

migrate-down-one:
	$(MIGRATE) down 1

run:
	go run cmd/api/main.go

gen:
	if not exist internal\httpapi\gen mkdir internal\httpapi\gen
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 -config oapi-codegen.yaml openapi.yaml

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint: lint-install
	"$(GOLANGCI_LINT)" run ./...
