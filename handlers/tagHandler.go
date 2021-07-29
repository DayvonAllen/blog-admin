package handlers

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TagHandler struct {
	TagService services.TagService
}

func (th *TagHandler) CreateTag(c *fiber.Ctx) error {
	c.Accepts("application/json")

	username := c.Locals("username").(string)

	tag := new(domain.Tag)

	err := c.BodyParser(tag)

	tag.Id = primitive.NewObjectID()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = th.TagService.Create(*tag, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (th *TagHandler) GetAllPostsByTags(c *fiber.Ctx) error {
	category := c.Params("category")
	page := c.Query("page", "1")

	postList, err := th.TagService.FindAllPostsByCategory(category, page)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
}

func (th *TagHandler) GetAllTags(c *fiber.Ctx) error {
	tags, err := th.TagService.FindAllTags()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": tags})
}
