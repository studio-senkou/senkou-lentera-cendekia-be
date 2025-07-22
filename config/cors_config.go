package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://portal.lenteracendekia.id,https://trialportal.lenteracendekia.id,http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	})
}
