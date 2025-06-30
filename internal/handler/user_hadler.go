package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"media-service/internal/services"
	"net/http"
)

type UserInput struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func LoginHandler(c *gin.Context) {
	var AuthRequest struct {
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"  binding:"required"`
		//rememberMe string `json:"rememberMe"`
	}

	if err := c.ShouldBindJSON(&AuthRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Printf("Error binding JSON: %v", err)
	}

	token, refresh_token, err := services.Authenticate(AuthRequest.Email, AuthRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refresh_token,
	})
	/*
		var maxAge int
		if loginRequest.RememberMe {
			maxAge = 30 * 24 * 60 * 60
		} else {
			maxAge = 0
		}

		c.SetCookie("jwt_token", token, maxAge, "/", domain, secure, httpOnly)
	*/
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
}
