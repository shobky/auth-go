package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `gorm:"not null"`
}

type RefreshToken struct {
	gorm.Model
	UserID uint   `gorm:"index"` // ID of the associated user
	Token  string `gorm:"not null"`
	Expiry int64  `gorm:"not null"`
}
