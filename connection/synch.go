package connection

import (
	"github.com/pashaakbar/blog-backend/models"
	"gorm.io/gorm"
)

func SyncDatabase(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
