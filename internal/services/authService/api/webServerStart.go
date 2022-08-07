package api

import (
	"messengerServer/internal/services/authService/api/handlers"
	"messengerServer/internal/services/authService/config"

	"github.com/gofiber/fiber/v2"
)

func WebStart(conf config.AuthServiceConfig) {
	app := fiber.New()
	app.Get("/authorize", handlers.AuthorizationHandler(conf))
	app.Get("/token", handlers.CheckTokenHandler(conf))
	app.Get("/register", handlers.RegistrationHandler(conf))
	app.Listen((":" + conf.Port))
}
