package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		fmt.Println("User role from context:", userRole, "Expected role:", role)
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Forbidden",
				"error":   "You do not have permission to access this resource",
			})
		}
		return c.Next()
	}
}
