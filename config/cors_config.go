package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://portal.lenteracendekia.id, https://trialportal.lenteracendekia.id, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	})
}
