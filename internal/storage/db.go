package storage

import (
	"log"

	"task-api/internal/domain"

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

	log.Println("db: migrate start")
	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		return nil, err
	}
	log.Println("db: migrate ok")

	return db, nil
}
