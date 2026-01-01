package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
	"github.com/studio-senkou/lentera-cendekia-be/utils/auth"
)

var jwtManager *auth.JwtManager

func init() {
	authSecret := app.GetEnv("AUTH_SECRET", "")
	jwtManager = auth.NewJwtManager(authSecret)
}

func AuthMiddleware() fiber.Handler {
	db := database.GetDB()
	authRepository := models.NewAuthenticationRepository(db)

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

		if err != nil || authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid authorization header format",
			})
		}

		claims, err := jwtManager.ValidateToken(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid token",
			})
		}

		// Safe type assertion untuk payload
		payload, exists := claims["payload"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid token payload",
			})
		}

		payloadMap, ok := payload.(map[string]any)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid token payload format",
			})
		}

		var userID int
		if id, ok := payloadMap["user_id"].(float64); ok {
			userID = int(id)
		} else {
			return c.Status(401).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid user ID in token",
				"error":   "Invalid user ID",
			})
		}

		userRole, ok := payloadMap["role"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid role in token",
				"error":   "Invalid role",
			})
		}

		if session, err := authRepository.GetTokenByUserID(userID); err != nil || session == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Session already invalidated or user logged out, please login again",
			})
		}

		c.Locals("userID", userID)
		c.Locals("userRole", userRole)

		return c.Next()
	}
}
