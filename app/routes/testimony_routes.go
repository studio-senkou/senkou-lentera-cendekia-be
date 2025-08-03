package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupTestimonyRoutes(router fiber.Router) {
	testimonyController := controllers.NewTestimonyController()

	router.Post("/testimonies", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), testimonyController.CreateTestimony)
	router.Get("/testimonies", testimonyController.GetAllTestimonials)
	router.Get("/testimonies/:id", testimonyController.GetTestimonyByID)
	router.Put("/testimonies/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), testimonyController.UpdateTestimony)
	router.Delete("/testimonies/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), testimonyController.DeleteTestimony)
}
