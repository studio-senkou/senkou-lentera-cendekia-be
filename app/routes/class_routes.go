package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupClassRoutes(router fiber.Router) {
	classController := controllers.NewClassController()

	router.Post("/classes", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), classController.CreateNewClass)
	router.Get("/classes", middlewares.AuthMiddleware(), classController.GetAllClasses)
	router.Get("/classes/dropdown", middlewares.AuthMiddleware(), classController.GetClassDropdown)
	router.Put("/classes/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), classController.UpdateClass)
	router.Delete("/classes/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), classController.DeleteClass)
}
