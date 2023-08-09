package database

import (
	"gofiber-auth/models"
)

func Sync() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.RefreshToken{})

}
