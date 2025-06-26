package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"media-service/internal/database"
	"media-service/internal/handler"
)

func main() {

	if err := database.InitDB(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	r := gin.Default()

	r.POST("/api/users", handler.CreateUser)

	log.Println("Starting server on :8000")

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
