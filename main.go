package main

import (
	"gofiber-auth/database"
	"gofiber-auth/handlers"
	"gofiber-auth/inits"
	"gofiber-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func init() {
	inits.LoadEnvs()
	database.Connect()
	database.Sync()
}
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Get("/validate", middleware.RequireAuth, handlers.Validate)
	app.Post("/refresh", handlers.RefreshToken)

	app.Listen(":3000")
}
