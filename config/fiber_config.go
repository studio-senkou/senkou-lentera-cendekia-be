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
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   message,
				"message": err.Error(),
			})
		},
	}
}
