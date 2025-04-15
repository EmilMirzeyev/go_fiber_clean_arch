package entity

import (
	"time"
)

type File struct {
	ID        uint      `gorm:"primaryKey"`
	FileName  string    `gorm:"size:255;not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
