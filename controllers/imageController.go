package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		fmt.Println(err)
		return err
	}

	files := form.File["image"]
	filename := ""

	for _, file := range files {
		filename = file.Filename
		fmt.Println(filename)

		err := c.SaveFile(file, "./uploads/"+filename)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return c.JSON(fiber.Map{
		"url": "http://localhost:9999/api/uploads/" + filename,
	})
}
