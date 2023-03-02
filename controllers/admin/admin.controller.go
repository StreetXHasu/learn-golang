package controllerAdmin

import (
	"app-learn-golang/initializers"
	"app-learn-golang/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func UserIndex(c *fiber.Ctx) error {
	//get all users from database
	var users []models.User
	result := initializers.DB.Find(&users)
	fmt.Println("result.RowsAffected", result.RowsAffected)
	fmt.Println("result.Error", result.Error)
	//result := resultRaw
	fmt.Println("users", users)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"users": users}})
}
