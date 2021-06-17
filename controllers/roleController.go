package controllers

import (
	"API/database"
	"API/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AllRoles(c *fiber.Ctx) error {
	var roles []models.Role

	database.DB.Preload("Permissions").Find(&roles)

	return c.JSON(roles)
}

type RoleDto struct {
	Id          uint
	Name        string
	Permissions []float64
}

func CreateRole(c *fiber.Ctx) error {
	var roleDto RoleDto

	err := c.BodyParser(&roleDto)

	if err != nil {
		return err
	}

	list := roleDto.Permissions

	permissions := make([]models.Permission, len(list))

	for i, permissionId := range list {
		permissions[i] = models.Permission{
			Id: uint(permissionId),
		}
	}

	role := models.Role{
		Name:        roleDto.Name,
		Permissions: permissions,
	}

	database.DB.Create(&role)

	return c.JSON(role)
}

func GetRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	database.DB.Preload("Permissions").Find(&role)

	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var roleDto RoleDto

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	list := roleDto.Permissions

	permissions := make([]models.Permission, len(list))

	for i, permissionId := range list {
		permissions[i] = models.Permission{
			Id: uint(permissionId),
		}
	}

	var result interface{}

	database.DB.Table("role_permissions").Where("role_id", id).Delete(&result)

	role := models.Role{
		Id:          uint(id),
		Name:        roleDto.Name,
		Permissions: permissions,
	}

	database.DB.Model(&role).Updates(role)

	return c.JSON(role)
}

func DeleteRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	database.DB.Delete(&role)

	return c.JSON(fiber.Map{
		"message": "Role successfully deleted",
	})
}
