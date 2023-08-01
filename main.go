package main

import (
	"os"
	"urbskali/file/routes"
	"urbskali/file/utils"

	"github.com/gofiber/fiber/v2"
)

func server() {
	utils.StartUp()

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024 * 1024, // 100 GB
	})

	routes.StaticRoutes(app)

	routes.PublicRoutes(app)
	routes.AdminRoutes(app)

	utils.Start(app)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		utils.Setup()
	} else if len(os.Args) > 1 && os.Args[1] == "build-ui" {
		utils.BuildUI()
	} else {
		server()
	}
}
