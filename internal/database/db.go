package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"media-service/internal/models"
	"media-service/pkg/config"
)

var GormDB *gorm.DB

func InitDB() error {
	config.GetDBConnString()

	log.Println(config.Cnfg.DBurl)
	var err error
	GormDB, err = gorm.Open(postgres.Open(config.Cnfg.DBurl), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB, err := GormDB.DB()

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to the database using gorm.")

	err = GormDB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return nil
}
