package handler

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/model"
	"media-service/internal/services"
	"net/http"
	"strconv"
)

func CreateProject(c *gin.Context) {
	var jsonBody model.ProjectInput

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	title := jsonBody.Title
	clientID := jsonBody.ClientID
	pmID := jsonBody.PmID

	id, err := services.CreateProject(title, clientID, pmID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetAllProjects(c *gin.Context) {
	projects, err := services.GetAllProjects()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func GetMyProjects(c *gin.Context) {
	emailVal, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in token"})
		return
	}

	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email format"})
		return
	}

	projects, err := services.GetMyProjects(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func GetProject(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID64, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	projectID := uint(projectID64)

	project, err := services.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}
