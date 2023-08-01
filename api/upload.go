package api

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	file_path := c.FormValue("path")
	fmt.Println("[UPLOAD] " + file_path)
	// Check is the secret is correct
	if c.FormValue("secret") != os.Getenv("OurTransfert_SECRET") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// URL decode the path
	// get the files from the request
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get the files",
		})
	}
	// get the files
	files := form.File["files"]
	// loop through the files
	for _, file := range files {
		// save the file
		fmt.Println("[Upload] Received: ", file.Filename)
		err := c.SaveFile(file, fmt.Sprintf("./files/%s/%s", file_path, file.Filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save the file",
			})
		}
	}
	return c.JSON(fiber.Map{
		"message": "Files uploaded successfully",
	})
}
