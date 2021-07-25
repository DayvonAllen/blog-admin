package handlers

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name: "Authentication",
		Value: string(signedToken),
		Expires: time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	c.Accepts("application/json")

	adminDetails := new(domain.Admin)

	err := c.BodyParser(adminDetails)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	conn, err := database.ConnectToDB()
	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

	if err != nil {
		panic(err)
	}

	admin := domain.Admin{Username: adminDetails.Username, Password: adminDetails.Password}
	admin.Id = primitive.NewObjectID()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	admin.Password = string(hashedPassword)

	_, err = conn.AdminCollection.InsertOne(context.TODO(), &admin)

	if err != nil {
		panic("error processing data")
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
