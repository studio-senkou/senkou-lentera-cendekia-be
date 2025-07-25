package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupUserRoutes(router fiber.Router) {
	userController := controllers.NewUserController()

	router.Post("/users", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.CreateStudent)
	router.Post("/users/mentors", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.CreateMentor)
	router.Post("/users/activate", userController.ActivateUser)
	router.Get("/users", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.GetAllUsers)
	router.Get("/users/students/dropdown", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.GetUserAsDropdown)
	router.Get("/users/mentors/dropdown", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.GetMentorDropdown)
	router.Get("/users/me", middlewares.AuthMiddleware()) // Get logged-in user
	router.Get("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.GetUser)
	router.Put("/users", middlewares.AuthMiddleware()) // Update logged-in user
	router.Put("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.UpdateUser)
	router.Delete("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), userController.DeleteUser)
}
