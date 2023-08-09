package initializers

import (
	"auth/models"
)

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
