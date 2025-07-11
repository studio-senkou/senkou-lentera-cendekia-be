package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName:                 "Lentera Quiz API",
		DisableStartupMessage:   true,
		CaseSensitive:           true,
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	}
}
