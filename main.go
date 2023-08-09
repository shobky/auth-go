package main

import (
	"auth/controllers"
	"auth/initializers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DbConnect()
	initializers.SyncDb()
}
func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Sinup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
