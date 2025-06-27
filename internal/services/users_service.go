package services

import (
	"errors"
	"media-service/internal/repository"
)

func CreateUser(email, password, role string) error {

	_, err := repository.GetUserByMail(email)

	if err == nil {
		return errors.New("user with this email already exists")
	}

	err = repository.CreateUser(email, password, role)
	if err != nil {
		return errors.New("could not create user in repository")
	}

	return nil
}
