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

func Refresh(refreshTokenString string) (string, string, error) {
	claims, err := jwt.ParseToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("could not parse refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid token type: expected refresh token")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("subject not found in token")
	}

	_, err = repository.GetUserByMail(email)
	if err != nil {
		return "", "", fmt.Errorf("user '%s' from token not found", email)
	}

	newAccessToken, newRefreshToken, err := jwt.GenerateTokens(email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate new tokens: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUser(email string, password string, roles []string) (*model.User, error) {

	_, err := repository.GetUserByMail(email)

	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := repository.CreateUser(email, passwordHash, roles)
	if err != nil {
		return nil, err
	}

	return user, nil
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
