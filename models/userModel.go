package models

import "gorm.io/gorm"

type User struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	Email    string `gorm:"unique"`
	Name     string
	Password string
	gorm.Model
}
