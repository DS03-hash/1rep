DB_DSN := postgres://postgres:postgres@localhost:5432/task_api?sslmode=disable
MIGRATE := migrate -path ./migrations -database "$(DB_DSN)"

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
