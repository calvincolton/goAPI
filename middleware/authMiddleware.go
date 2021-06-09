package middleware

import (
	"API/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	if _, err := utils.ParseJWT(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)
		fmt.Println(err)
		return c.JSON(fiber.Map{
			"message": "You have not been authenticated",
		})
	}

	return c.Next()
}