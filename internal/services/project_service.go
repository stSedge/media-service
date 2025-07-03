package services

import (
	"errors"
	"media-service/internal/model"
	"media-service/internal/repository"
)

func CreateProject(title string, clientID uint, pmID uint) (uint, error) {
	id, err := repository.CreateProject(title, clientID, pmID)

	if err != nil {

		return 0, errors.New("could not create project in repository")
	}

	return id, nil
}

func GetAllProjects() ([]model.ProjectResponse, error) {
	projects, err := repository.GetAllProjects()
	if err != nil {
		return nil, errors.New("could not get projects")
	}

	var response []model.ProjectResponse
	for _, p := range projects {
		response = append(response, model.ProjectResponse{
			ID:    p.ID,
			Title: p.Title,
			Client: model.UserResponse{
				ID:    p.Client.ID,
				Email: p.Client.Email,
			},
			Pm: model.UserResponse{
				ID:    p.Pm.ID,
				Email: p.Pm.Email,
			},
		})
	}

	return response, nil
}

func GetMyProjects(email string) ([]model.ProjectResponse, error) {
	user, ok := repository.GetUserByMail(email)
	if ok != nil {
		return nil, errors.New("could not get user")
	}

	projects, err := repository.GetMyProjects(user.ID)
	if err != nil {

		return nil, errors.New("could not get projects")
	}

	var response []model.ProjectResponse
	for _, p := range projects {
		response = append(response, model.ProjectResponse{
			ID:    p.ID,
			Title: p.Title,
			Client: model.UserResponse{
				ID:    p.Client.ID,
				Email: p.Client.Email,
			},
			Pm: model.UserResponse{
				ID:    p.Pm.ID,
				Email: p.Pm.Email,
			},
		})
	}

	return response, nil
}
