package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/controllers"
	"github.com/studio-senkou/lentera-cendekia-be/app/middlewares"
)

func SetupStaticAssetRoutes(router fiber.Router) {
	staticAssetController := controllers.NewStaticAssetController()

	router.Post("/static-assets", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), staticAssetController.CreateStaticAsset)
	router.Get("/static-assets", staticAssetController.GetAllStaticAssets)
	router.Get("/static-assets/:id", staticAssetController.GetStaticAssetByID)
	router.Delete("/static-assets/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), staticAssetController.DeleteStaticAsset)
}