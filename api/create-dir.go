package api

import (
	"fmt"
	"os"
	"urbskali/file/state"

	"github.com/gofiber/fiber/v2"
)

func CreateDir(c *fiber.Ctx) error {
	if c.FormValue("secret") != state.Config.Secret {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// URL decode the path
	path := c.FormValue("path")
	name := c.FormValue("name")
	// create the directory
	fmt.Println("[CREATE] " + path + name)
	err := os.MkdirAll(fmt.Sprintf("./files/%s/%s", path, name), 0755)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create the directory",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Directory created successfully",
	})
}
