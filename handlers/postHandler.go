package handlers

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type PostHandler struct {
	PostService services.PostService
}

func (ph *PostHandler) CreatePost(c *fiber.Ctx) error {
	c.Accepts("application/json")

	username := c.Locals("username").(string)

	post := new(domain.Post)

	err := c.BodyParser(post)

	post.Author = username
	post.Id = primitive.NewObjectID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ph.PostService.Create(*post, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PostHandler) UpdatePost(c *fiber.Ctx) error {
	c.Accepts("application/json")

	username := c.Locals("username").(string)

	post := new(domain.PostUpdateDto)

	err := c.BodyParser(post)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ph.PostService.UpdateByTitle(*post, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PostHandler) UpdateVisibility(c *fiber.Ctx) error {
	c.Accepts("application/json")

	username := c.Locals("username").(string)

	post := new(domain.PostUpdateVisibilityDto)

	err := c.BodyParser(post)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ph.PostService.UpdateVisibility(*post, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newStoriesQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newStoriesQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	postList, err := ph.PostService.FindAllPosts(page, isNew)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
}

func (ph *PostHandler) GetPostById(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	post, err := ph.PostService.FindPostById(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": post})
}