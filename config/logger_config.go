package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewLoggerConfig() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${status} ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Output:     nil,
	})
}