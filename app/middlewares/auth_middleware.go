package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
)

var jwtManager *auth.JwtManager

func init() {
	authSecret := app.GetEnv("AUTH_SECRET", "")
	jwtManager = auth.NewJwtManager(authSecret)
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims["payload"])

		return c.Next()
	}
}
