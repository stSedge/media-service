package model

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	JTI       uuid.UUID `gorm:"type:uuid;unique;not null;index"`
	UserID    uint      `gorm:"not null;index"`
	IsActive  bool      `gorm:"not null;default:true"`
	SessionID uint      `gorm:"not null"`
}

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiredAt time.Time `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	IPAddress string    `gorm:"not null"`
	Expired   bool      `gorm:"not null;default:false"`
}

type SessionResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	Expired   bool      `json:"expired"`
	IsCurrent bool      `json:"is_current"`
}
