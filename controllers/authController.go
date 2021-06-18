package controllers

import (
	"API/database"
	"API/models"
	"API/utils"

	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    3,
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	token, err := utils.GenerateJWT(strconv.Itoa(int(user.Id)))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     utils.CookieSecret,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successfully logged in",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     utils.CookieSecret,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies(utils.CookieSecret)

	id, _ := utils.ParseJWT(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	cookie := c.Cookies(utils.CookieSecret)

	id, _ := utils.ParseJWT(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)

	if err != nil {
		return err
	}
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	cookie := c.Cookies(utils.CookieSecret)
	id, _ := utils.ParseJWT(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id: uint(userId),
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}
