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
		userGroup := protect.Group("/")
		userGroup.Use(middleware.CheckRole("ROLE_USER"))
		{
			userGroup.POST("/logout", handler.LogoutHandler)
			userGroup.POST("/logout/all", handler.LogoutAllHandler)
			userGroup.POST("/refresh", handler.RefreshTokenHandler)
		}

		pmGroup := protect.Group("/")
		pmGroup.Use(middleware.CheckRole("ROLE_PM"))
		{
			pmGroup.GET("/projects/my", handler.GetMyProjects)
			pmGroup.GET("/projects/:project_id", handler.GetProject)
		}

		adminGroup := protect.Group("/")
		adminGroup.Use(middleware.CheckRole("ROLE_ADMIN"))
		{
			adminGroup.POST("/users", handler.CreateUser)
			adminGroup.GET("/users", handler.GetAllUsers)
			adminGroup.POST("/projects", handler.CreateProject)
			adminGroup.GET("/projects", handler.GetAllProjects)
			adminGroup.POST("/projects/:project_id/reports", handler.CreateReport)
		}
	}

	log.Println("Starting server on :8000")

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
