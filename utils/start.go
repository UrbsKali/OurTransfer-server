package utils

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Start(app *fiber.App) {
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if os.Getenv("HTTPS") == "true" {
		app.ListenTLS(port, os.Getenv("CERT"), os.Getenv("KEY"))
	} else {
		app.Listen(port)
	}
}
