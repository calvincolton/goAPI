package controllers

import (
	"API/database"
	"API/models"

	"github.com/gofiber/fiber/v2"
)

func AllPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission
	database.DB.Find(&permissions)

	return c.JSON(permissions)
}

func CreatePermission(c *fiber.Ctx) error {
	var permission models.Permission

	err := c.BodyParser(&permission)

	if err != nil {
		return err
	}

	database.DB.Create(&permission)

	return c.JSON(permission)
}
