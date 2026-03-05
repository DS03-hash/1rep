package storage

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(dsn string) (*gorm.DB, error) {
	log.Println("db: opening connection...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	log.Println("db: ping ok")

	var currentDB string
	if err := db.Raw("SELECT current_database()").Scan(&currentDB).Error; err != nil {
		return nil, err
	}
	log.Println("db: current_database =", currentDB)

	return db, nil
}
// OpenDB - это функция, которая открывает соединение с базой данных PostgreSQL с помощью GORM.
// Она принимает строку подключения (DSN) в качестве аргумента и возвращает указатель на gorm.DB и ошибку, если операция не удалась.
// Внутри функции выполняются следующие шаги:
// 1. Открытие соединения с базой данных с помощью gorm.Open и драйвера postgres.
// 2. Получение объекта sql.DB из gorm.DB для выполнения операции Ping, чтобы проверить, что соединение успешно установлено.
// 3. Выполнение SQL-запроса для получения текущей базы данных и вывод ее имени в лог.
// 4. Возвращение указателя на gorm.DB и nil, если все операции прошли успешно, или возвращение ошибки, если что-то пошло не так.
