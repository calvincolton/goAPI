package controllers

import (
	"API/database"
	"API/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AllOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB, &models.Order{}, page))
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	err := c.BodyParser(&order)

	if err != nil {
		return err
	}

	database.DB.Create(&order)

	return c.JSON(order)
}
