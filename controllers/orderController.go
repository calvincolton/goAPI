package controllers

import (
	"API/database"
	"API/models"
	"encoding/csv"
	"os"
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

func ExportOrder(c *fiber.Ctx) error {
	filePath := "./csv/orders.csv"

	err := CreateFile(filePath)
	if err != nil {
		return err
	}

	return c.Download(filePath)
}

func CreateFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	writer.Write([]string{
		"ID",
		"Name",
		"Email",
		"Product Title",
		"Price",
		"Quantity",
	})

	for _, order := range orders {
		data := []string{
			strconv.Itoa(int(order.Id)),
			order.FirstName + " " + order.LastName,
			order.Email,
			"",
			"",
			"",
		}

		err := writer.Write(data)
		if err != nil {
			return err
		}

		for _, orderItem := range order.OrderItems {
			data := []string{
				"",
				"",
				"",
				orderItem.ProductTitle,
				strconv.Itoa(int(orderItem.Price)),
				strconv.Itoa(int(orderItem.Quantity)),
			}

			err := writer.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

func Chart(c *fiber.Ctx) error {
	var sales []Sales

	database.DB.Raw(`
		SELECT DATE_FORMAT(o.created_at, '%Y-%m-%d') as date, 
		SUM(oi.price * oi.quantity) as sum 
		FROM orders o 
		JOIN order_items oi ON o.id = oi.order_id 
		GROUP BY date;
	`).Scan(&sales)

	return c.JSON(sales)
}
