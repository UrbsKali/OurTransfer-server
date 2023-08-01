package api

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CheckSecret(c *fiber.Ctx) error {
	fmt.Println("[Check Secret] Checking secret from IP:", c.IP())
	// URL decode the path
	input_secret := c.FormValue("secret")
	if input_secret == os.Getenv("OurTransfert_SECRET") {
		return c.JSON(fiber.Map{
			"message": true,
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Incorrect secret",
	})
}
