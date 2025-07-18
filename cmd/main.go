package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"media-service/internal/database"
	"media-service/internal/handler"
	"media-service/internal/middleware"
)

func main() {

	if err := database.InitDB(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	r := gin.Default()

	r.POST("/api/login", handler.LoginHandler)

	protect := r.Group("/api")

	protect.Use(middleware.JWTMiddleware())
	{
		protect.POST("/logout", handler.LogoutHandler)
		protect.POST("/logout/all", handler.LogoutAllHandler)
		protect.GET("/sessions", handler.GetSessionsHandler)
		r.POST("/api/users", handler.CreateUser)
		r.POST("/refresh", handler.RefreshTokenHandler)
		r.GET("/api/users", handler.GetAllUsers)
		r.POST("/api/projects", handler.CreateProject)
		r.GET("/api/projects", handler.GetAllProjects)
		r.GET("/api/projects/my", handler.GetMyProjects)
		r.GET("/api/projects/:project_id", handler.GetProject)
		r.POST("/api/projects/:project_id/reports", handler.CreateReport)
	}

	log.Println("Starting server on :8000")

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
