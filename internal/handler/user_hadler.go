package handler

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/services"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var jsonBody map[string]string

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	email := jsonBody["email"]
	password := jsonBody["password"]
	role := jsonBody["role"]

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	err := services.CreateUser(email, password, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ...
}
