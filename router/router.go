package router

import (
	"com.aharakitchen/app/handlers"
	"com.aharakitchen/app/repo"
	"com.aharakitchen/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	th := handlers.TagHandler{TagService: services.NewTagService(repo.NewTagRepoImpl())}
	ph := handlers.PostHandler{PostService: services.NewPostService(repo.NewPostRepoImpl())}
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	tags := api.Group("/control/tags")
	tags.Post("/", th.CreateTag)

	posts := api.Group("/control/posts")
	posts.Post("/", ph.CreatePost)
	posts.Put("/visibility", ph.UpdateVisibility)
	posts.Put("/", ph.UpdatePost)

	auth := api.Group("/control/checkin")
	auth.Post("/", ah.Login)
}

func Setup() *fiber.App {
	app := fiber.New()

	SetupRoutes(app)
	return app
}
