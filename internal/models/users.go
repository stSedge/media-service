package models

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Password     string `json:"password" db:"-"`
	PasswordHash string `json:"-" db:"password_hash"`
	// Roles     []string `json:"roles" db:"roles"`
}
