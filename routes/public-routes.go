package routes

import (
	"urbskali/file/api"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App) {

	route := app.Group("/api")

	route.Post("/get_files/", api.GetFiles)
	route.Post("/file_info/", api.FileInfo)
	route.Post("/download/", api.Download)
	route.Post("/check_secret/", api.CheckSecret)

}
