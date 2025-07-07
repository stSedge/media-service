package repository

import (
	"log"
	"media-service/internal/database"
	"media-service/internal/model"
)

func CreateReport(projectID uint, title string, content string, members string, filename string) (uint, error) {
	report := &model.Report{
		Title:     title,
		ProjectID: projectID,
		Content:   content,
		FilePath:  filename,
		Members:   members,
	}
	res := database.GormDB.Create(&report)

	if res.Error != nil {
		log.Printf("Error creating report: %v", res.Error)
		return 0, res.Error
	}
	log.Printf("Report created")
	return report.ID, nil
}
