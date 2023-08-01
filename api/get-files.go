package api

import (
	"urbskali/file/utils"

	"github.com/gofiber/fiber/v2"
)

func GetFiles(c *fiber.Ctx) error {
	files, err := utils.GetFiles(c.FormValue("path"))
	if err != nil {
		return c.Status(404).SendString("Directory not found")
	}
	return c.JSON(files)
}
