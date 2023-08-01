package utils

import (
	"fmt"
	"urbskali/file/state"

	"github.com/gofiber/fiber/v2"
)

func Start(app *fiber.App) {
	port := fmt.Sprintf(":%s", state.Config.Port)
	if state.Config.HTTPS {
		app.ListenTLS(port, state.Config.Cert, state.Config.Key)
	} else {
		app.Listen(port)
	}
}
