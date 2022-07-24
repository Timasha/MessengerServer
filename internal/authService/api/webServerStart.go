package api

import (
	"logger/internal/authService/api/handlers"
	"logger/internal/authService/config"

	"github.com/gofiber/fiber/v2"
)

func WebStart(conf config.AuthServiceConfig) {
	app := fiber.New()
	app.Get("/authorize", handlers.AuthorizeHandler(conf))
	app.Listen((":" + conf.Port))
}
