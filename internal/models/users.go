package models

type User struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	Email        string   `json:"email" gorm:"unique;not null"`
	PasswordHash string   `json:"-" gorm:"not null"`
	Roles        []string `json:"roles" gorm:"type:text[]"`
}
