package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
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
