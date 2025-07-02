package model

type User struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	Email        string   `json:"email" gorm:"unique;not null"`
	PasswordHash string   `json:"-" gorm:"not null"`
	Roles        []string `json:"roles" gorm:"type:text[]"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
