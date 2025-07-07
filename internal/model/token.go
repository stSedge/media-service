package model

import (
	"github.com/google/uuid"
)

type Token struct {
	ID       uint      `gorm:"primaryKey"`
	JTI      uuid.UUID `gorm:"type:uuid;unique;not null;index"`
	UserID   uint      `gorm:"not null;index"`
	IsActive bool      `gorm:"not null;default:true"`
}


