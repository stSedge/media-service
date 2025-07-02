package model

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Roles        pq.StringArray `json:"roles" gorm:"type:text[]"`
	CreatedAt    time.Time      `json:"-" gorm:"autoCreateTime"`
}

type CreatedAtInfo struct {
	Date         time.Time `json:"date"`
	TimezoneType int       `json:"timezone_type"`
	Timezone     string    `json:"timezone"`
}

type UserResponseFull struct {
	ID        uint          `json:"id"`
	Email     string        `json:"email"`
	Roles     []string      `json:"roles"`
	CreatedAt CreatedAtInfo `json:"createdAt"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type UserInput struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}
