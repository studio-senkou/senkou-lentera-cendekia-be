package middlewares

import (
	"fmt"

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
				"status":  "fail",
				"message": "Unauthorized",
				"error":   "Missing or invalid token",
			})
		}

		var authToken string

		_, err := fmt.Sscanf(token, "Bearer %s", &authToken)

		// If the token is not in the correct format, we return an Uauthorized response
		// The correct format is "Bearer <token>"
		if err != nil || authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid authorization header format",
			})
		}

		// Validat the token using the JWT manager
		// If the token is invalid, we return an Unauthorized response
		// The ValidateToken method will return the claims if the token is valid
		claims, err := jwtManager.ValidateToken(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid token",
			})
		}

		c.Locals("userID", claims["payload"])

		return c.Next()
	}
}
