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

	// Create user in database
	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create user",
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
			"error": "Failed to read body",
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
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create token ",
		})
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Refresh token expiration time
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create refresh token",
		})
	}

	// Store refresh token in database
	refreshTokenRecord := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshTokenString,
		Expiry: time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	result := database.DB.Create(&refreshTokenRecord)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to store refresh token",
		})
	}

	// Set cookies for access and refresh tokens
	accessCookie := &fiber.Cookie{
		Name:     "AccessToken",
		Value:    accessTokenString,
		Expires:  time.Now().Add(time.Hour * 2),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(accessCookie)

	return c.SendStatus(http.StatusOK)
}

func RefreshToken(c *fiber.Ctx) error {
	// Get user ID from request body
	var body struct {
		UserID uint
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// Fetch refresh token from database using userID
	var refreshTokenRecord models.RefreshToken
	result := database.DB.First(&refreshTokenRecord, "user_id = ?", body.UserID)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to fetch refresh token from database",
		})
	}

	// Verify and parse the refresh token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenRecord.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	// Get the user ID from the claims
	userID, ok := claims["sub"].(float64)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid refresh token format",
		})
	}

	// Look up the user in the database
	var user models.User
	result = database.DB.First(&user, uint(userID))
	if result.Error != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Generate a new access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create access token",
		})
	}

	// Set the new access token in a cookie
	accessCookie := &fiber.Cookie{
		Name:     "AccessToken",
		Value:    accessTokenString,
		Expires:  time.Now().Add(time.Hour * 2),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(accessCookie)

	return c.SendStatus(http.StatusOK)
}

func Validate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get user from context",
		})
	}

	return c.JSON(fiber.Map{
		"message": user,
	})
}
