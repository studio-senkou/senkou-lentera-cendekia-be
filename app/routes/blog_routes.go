package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupBlogRoutes(router fiber.Router) {
	blogController := controllers.NewBlogController()

	router.Post("/blogs", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("mentor"), blogController.CreateBlog)
	router.Get("/blogs", blogController.GetAllBlogs)
	router.Get("/blogs/:id", blogController.GetBlogByID)
	router.Put("/blogs/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("mentor"), blogController.UpdateBlog)
	router.Delete("/blogs/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("mentor"), blogController.DeleteBlog)
}
