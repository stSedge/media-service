package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"media-service/internal/model"
	"media-service/internal/repository"
	"media-service/pkg/jwt"
)

func Authenticate(email, password string) (string, string, error) {
	user, err := repository.GetUserByMail(email)

	if err != nil {
		log.Printf("Error finding user by email %s: %v", email, err)
		return "", "", err // Возвращаем ошибку дальше
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	return jwt.GenerateTokens(user.Email)
}

func CreateUser(email string, password string, roles []string) error {

	_, err := repository.GetUserByMail(email)

	if err == nil {
		return errors.New("user with this email already exists")
	}

	err = repository.CreateUser(email, password, roles)
	if err != nil {
		return errors.New("could not create user in repository")
	}

	return nil
}

func GetAllUsers() ([]model.UserResponseFull, error) {
	users, err := repository.GetAllUsers()

	if err != nil {
		return nil, errors.New("could not get all users")
	}

	var usersResponse []model.UserResponseFull
	for _, u := range users {
		tzName, _ := u.CreatedAt.Zone()

		usersResponse = append(usersResponse, model.UserResponseFull{
			ID:    u.ID,
			Email: u.Email,
			Roles: u.Roles,
			CreatedAt: model.CreatedAtInfo{
				Date:         u.CreatedAt,
				TimezoneType: 3, // заглушка! в го нет аналога time_zone из РЗ
				Timezone:     tzName,
			},
		})
	}

	return usersResponse, nil
}
