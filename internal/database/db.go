package database

import (
	"fmt"
	"log"
	"media-service/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() error {
	config.GetDBConnString()

	log.Println(config.Cnfg.DBurl)
	var err error
	DB, err = sqlx.Connect("postgres", config.Cnfg.DBurl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close()
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to the database using sqlx.")
	return nil
}
