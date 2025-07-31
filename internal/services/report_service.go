package services

import (
	"errors"
	"media-service/internal/repository"
)

func CreateReport(projectID uint, title string, content string, members string, filename string) (uint, error) {
	id, err := repository.CreateReport(projectID, title, content, members, filename)
	if err != nil {
		return 0, errors.New("could not create project in repository")
	}
	return id, nil
}
