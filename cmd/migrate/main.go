package main

import (
	"log"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Conversation{},
		&models.Message{},
	); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("Migration completed successfully")
}
