package entity

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	Age       int       `gorm:"not null"`
	ImageName string    `gorm:"size:255"`
	RoleID    uint      `gorm:"not null"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
