package config

import (
	"fmt"
	"log"
	"ppob-backend/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *SystemConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPass, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err.Error())
	}

	db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Product{}, &model.Transaction{})

	return db
}
