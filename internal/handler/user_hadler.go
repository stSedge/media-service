package handler

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/services"
	"net/http"
)

type UserInput struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func CreateUser(c *gin.Context) {
	var jsonBody UserInput

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	email := jsonBody.Email
	password := jsonBody.Password
	roles := jsonBody.Roles

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	err := services.CreateUser(email, password, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ...
}
