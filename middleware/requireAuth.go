package middleware

import (
	"fmt"
	"gofiber-auth/database"
	"gofiber-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *fiber.Ctx) error {
	// Get the cookie
	tokenString := c.Cookies("AccessToken")
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}
	// Docore and validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	// Check the exp
	exp, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > exp {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}
	// Find the user with the token
	var user models.User
	database.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}
	// Attach to req
	c.Locals("user", user)

	// Continue
	return c.Next()
}
