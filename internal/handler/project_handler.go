package handler

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/model"
	"media-service/internal/services"
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{"data": id})
}

func GetAllProjects(c *gin.Context) {
	projects, err := services.GetAllProjects()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}
