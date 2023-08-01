package api

import (
	"fmt"
	"urbskali/file/state"

	"github.com/gofiber/fiber/v2"

	"os"
)

func Delete(c *fiber.Ctx) error {
	path := c.FormValue("path")
	if c.FormValue("secret") != state.Config.Secret {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	fmt.Println("[DELETE] " + path)
	// delete the file
	err := os.RemoveAll(fmt.Sprintf("./files/%s", path))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete the file",
		})
	}
	return c.JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
