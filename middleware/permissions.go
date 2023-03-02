package middleware

import (
	"app-learn-golang/models"
	"github.com/gofiber/fiber/v2"
)

func Permissions(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You do not have permission to access this resource"})
	}

	return c.Next()
}
