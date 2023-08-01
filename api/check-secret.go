package api

import (
	"fmt"
	"urbskali/file/state"

	"github.com/gofiber/fiber/v2"
)

func CheckSecret(c *fiber.Ctx) error {
	fmt.Println("[Check Secret] Checking secret from IP:", c.IP())
	input_secret := c.FormValue("secret")
	if input_secret == state.Config.Secret {
		return c.JSON(fiber.Map{
			"message": true,
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Incorrect secret",
	})
}
