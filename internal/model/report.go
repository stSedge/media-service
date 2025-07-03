package model

import "time"

type Report struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Title     string    `json:"title"      gorm:"not null"`
	Content   string    `json:"content"    gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	FilePath  string    `json:"file_path"  gorm:"not null"`
	ProjectID uint      `json:"-"          gorm:"not null; index"`
}
