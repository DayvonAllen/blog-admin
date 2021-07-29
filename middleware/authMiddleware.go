package middleware

import (
	"com.aharakitchen/app/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func IsLoggedIn(c *fiber.Ctx) error {
	token := c.Cookies("Authentication")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	c.Locals("username", u.Username)

	err = c.Next()

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	return nil
}
