package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
)

func SetupUserRoutes(router fiber.Router) {
	userController := controllers.NewUserController()

	router.Post("/users", userController.CreateStudent)
	router.Post("/users/mentors", userController.CreateMentor)
	router.Post("/users/activate", userController.ActivateUser)
	router.Get("/users", userController.GetAllUsers)
	router.Get("/users/students/dropdown", userController.GetUserAsDropdown)
	router.Get("/users/mentors/dropdown", userController.GetMentorDropdown)
	router.Get("/users/:id", userController.GetUser)
	router.Put("/users/:id", userController.UpdateUser)
	router.Delete("/users/:id", userController.DeleteUser)
}
