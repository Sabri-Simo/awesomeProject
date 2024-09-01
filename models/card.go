package models

import (
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Books  []Book `gorm:"many2many:card_books;" json:"books"`
}
