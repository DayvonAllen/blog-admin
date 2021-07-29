package router

import (
	"com.aharakitchen/app/handlers"
	"com.aharakitchen/app/middleware"
	"com.aharakitchen/app/repo"
	"com.aharakitchen/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	tags.Post("/", middleware.IsLoggedIn, th.CreateTag)
	tags.Get("/:category", middleware.IsLoggedIn, th.GetAllPostsByTags)
	tags.Get("/", middleware.IsLoggedIn, th.GetAllTags)

	posts := api.Group("/control/posts")
	posts.Post("/", middleware.IsLoggedIn, ph.CreatePost)
	posts.Put("/visibility", middleware.IsLoggedIn, ph.UpdateVisibility)
	posts.Put("/", middleware.IsLoggedIn, ph.UpdatePost)
	posts.Get("/:id", middleware.IsLoggedIn, ph.GetPostById)
	posts.Get("/", middleware.IsLoggedIn, ph.GetAllPosts)

	auth := api.Group("/control/checkin")
	auth.Get("/status", middleware.IsLoggedIn, ah.IsLoggedIn)
	auth.Post("/", ah.Login)
	auth.Get("/", ah.Logout)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)
	return app
}
