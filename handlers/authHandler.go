package handlers

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	c.Accepts("application/json")
	details := new(domain.LoginDetails)
	err := c.BodyParser(details)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	var auth domain.Authentication

	_, token, err := ah.AuthService.Login(strings.ToLower(details.Username), details.Password, c.IP(), c.IPs())

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Authentication failure")})
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken([]byte(token))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	signedToken = append(signedToken, t...)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": string(signedToken)})
}

func (ah *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "Authentication",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ah *AuthHandler) IsLoggedIn(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": c.Locals("username")})
}
