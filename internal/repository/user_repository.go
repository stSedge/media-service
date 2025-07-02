package repository

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"media-service/internal/database"
	"media-service/internal/model"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUser(email string, password string, roles []string) error {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: passwordHash,
		Roles:        roles,
	}
	res := database.GormDB.Create(&user)

	if res.Error != nil {
		log.Printf("Error creating user: %v", res.Error)
		return res.Error
	}

	log.Printf("User created")

	return nil
}

func GetUserByMail(email string) (*model.User, error) {
	var user model.User
	err := database.GormDB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID int) (*model.User, error) {
	var user model.User
	err := database.GormDB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User

	err := database.GormDB.Find(&users).Error

	if err != nil {
		log.Printf("Error fetching users: %v", err.Error)
		return nil, err
	}

	return users, nil
}
