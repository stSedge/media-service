package repository

import (
	"log"
	"media-service/internal/database"
	"media-service/internal/model"
)

func CreateProject(title string, clientID uint, pmID uint) (uint, error) {
	project := &model.Project{
		Title:    title,
		ClientID: clientID,
		PmID:     pmID,
	}
	res := database.GormDB.Create(&project)
	if res.Error != nil {
		log.Printf("Error creating project: %v", res.Error)
		return 0, res.Error
	}
	log.Printf("Project created")
	return project.ID, nil
}

func GetAllProjects() ([]model.Project, error) {
	var projects []model.Project

	res := database.GormDB.Preload("Client").Preload("Pm").Find(&projects)
	if res.Error != nil {
		log.Printf("Error fetching project: %v", res.Error)
		return nil, res.Error
	}
	log.Printf("Projects fetched")
	return projects, nil
}

func GetMyProjects(pmID uint) ([]model.Project, error) {
	var projects []model.Project

	res := database.GormDB.
		Preload("Client").
		Preload("Pm").
		Where("pm_id = ?", pmID).
		Find(&projects)

	if res.Error != nil {
		log.Printf("Error fetching project: %v", res.Error)
		return nil, res.Error
	}
	log.Printf("Projects fetched")
	return projects, nil
}
