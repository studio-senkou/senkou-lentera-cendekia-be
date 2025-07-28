package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupAuthRoutes(router fiber.Router) {
	authController := controllers.NewAuthController()

	router.Post("/auth/login", authController.LoginUser)
	router.Post("/auth/login/admin", authController.LoginAdmin)
	router.Post("/auth/verify-email", authController.VerifyAccount)
	router.Post("/auth/verify-token", authController.VerifyOneTimeToken)
	router.Put("/auth/refresh", authController.RefreshToken)
	router.Delete("/auth/logout", middlewares.AuthMiddleware(), authController.Logout)
}
