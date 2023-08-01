package api

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CreateDir(c *fiber.Ctx) error {
	if c.FormValue("secret") != os.Getenv("OurTransfert_SECRET") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// URL decode the path
	file, _ := url.QueryUnescape(c.Params("*"))
	// create the directory
	fmt.Println("[CREATE] " + file + c.FormValue("name"))
	err := os.MkdirAll(fmt.Sprintf("./files/%s/%s", file, c.FormValue("name")), 0755)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create the directory",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Directory created successfully",
	})
}
