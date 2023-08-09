package handlers

import (
	"gofiber-auth/database"
	"gofiber-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	// Get req body
	var body struct {
		Email    string
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// Has password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Respond
	return c.JSON(http.StatusOK)
}

func Login(c *fiber.Ctx) error {
	// Get req body
	var body struct {
		Email    string
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Look up user in db
	var user models.User
	database.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Email ",
		})
	}

	// Check passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Password ",
		})
	}
	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create token ",
		})
	}
	// Set cookie
	cookie := &fiber.Cookie{
		Name:     "Authentication",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	}

	c.Cookie(cookie)

	return c.SendStatus(http.StatusOK)
}

func Validate(c *fiber.Ctx) error {
	user, exists := c.Locals("user").(string)

	if exists {
		return c.JSON(fiber.Map{
			"message": user,
		})
	}
	return c.JSON(fiber.Map{
		"message": "Logged in, no user received",
	})
}
