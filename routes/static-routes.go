package routes

import (
	"github.com/gofiber/fiber/v2"
)

func StaticRoutes(app *fiber.App) {

	app.Static("/", "./ui")

	app.Get("/download/*", func(c *fiber.Ctx) error {
		return c.SendFile("./ui/index.html")
	})

}
