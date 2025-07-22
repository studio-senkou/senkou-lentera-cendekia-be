package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupAuthRoutes(router fiber.Router) {
	authController := controllers.NewAuthController()

	router.Post("/auth/login", authController.Login)
	router.Put("/auth/refresh", authController.RefreshToken)
	router.Delete("/auth/logout", middlewares.AuthMiddleware(), authController.Logout)
}
