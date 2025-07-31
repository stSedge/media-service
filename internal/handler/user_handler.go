package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"media-service/internal/model"
	"media-service/internal/repository"
	"media-service/internal/services"
	"net/http"
	"strings"
)

func LoginHandler(c *gin.Context) {
	var AuthRequest struct {
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"  binding:"required"`
		//rememberMe string `json:"rememberMe"`
	}

	if err := c.ShouldBindJSON(&AuthRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Printf("Error binding JSON: %v", err)
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	token, refresh_token, err := services.Authenticate(AuthRequest.Email, AuthRequest.Password, ipAddress, userAgent)
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

func LogoutHandler(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Refresh-Token header is required"})
		return
	}

	if err := services.Logout(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func LogoutAllHandler(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, refresh_token is required"})
		return
	}

	email := c.GetString("user_email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in token"})
		return
	}

	if err := services.LogoutAll(email, request.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out from all devices"})
}

func GetSessionsHandler(c *gin.Context) {
	email, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in token"})
		return
	}

	user, err := repository.GetUserByMail(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user"})
		return
	}

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Refresh-Token header is required"})
		return
	}

	sessions, err := services.GetUserSessions(user.ID, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func RefreshTokenHandler(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Refresh-Token header is required"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header is missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header format"})
		return
	}

	newAccessToken, newRefreshToken, err := services.Refresh(refreshToken)
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func CreateUser(c *gin.Context) {
	var jsonBody model.UserInput

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

	user, err := services.CreateUser(email, password, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"roles": user.Roles,
		},
	})
}

func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
