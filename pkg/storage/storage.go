package storage

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"user_crud/internal/domain/entity"
)

func NewDatabaseConnection(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.User{}, &entity.File{}, &entity.Role{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
