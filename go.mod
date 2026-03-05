module task-api

go 1.25.7

require (
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.5
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)
// В этом файле go.mod определены зависимости проекта и их версии.
// Он указывает, что проект использует Go версии 1.25.7 и требует
// две основные зависимости: gorm.io/driver/postgres для работы с базой данных PostgreSQL и gorm.io/gorm для использования ORM GORM.
// Также указаны дополнительные зависимости, которые являются косвенными (indirect) и используются внутри
// основных зависимостей. Эти зависимости включают различные библиотеки для работы с базой данных, тестирования и других утилит, которые необходимы для функционирования проекта.