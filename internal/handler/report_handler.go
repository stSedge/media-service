package handler

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/services"
	"net/http"
	"strconv"
)

func CreateReport(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID64, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}
	projectID := uint(projectID64)
	title := c.PostForm("title")
	content := c.PostForm("content")
	members := c.PostForm("members")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// TODO: сохранять файл в S3

	id, err := services.CreateReport(projectID, title, content, members, file.Filename)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
