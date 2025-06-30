package repository

import (
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"media-service/internal/database"
	"media-service/internal/models"
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
	query := `INSERT INTO users (email, password_hash, roles) VALUES ($1, $2, $3::user_role[])`

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	res, err := database.DB.Exec(query, email, passwordHash, pq.Array(roles))
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	log.Printf("User created. Rows affected: %d", rowsAffected)

	return nil
}

func GetUserByMail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, roles FROM users WHERE email=$1`
	var user models.User
	err := database.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID int) (*models.User, error) {
	query := `SELECT id, email, password_hash, roles FROM users WHERE id=$1`
	var user models.User
	err := database.DB.Get(&user, query, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
