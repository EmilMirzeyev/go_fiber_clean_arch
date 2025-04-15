package entity

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:50;not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
