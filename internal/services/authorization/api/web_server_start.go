package api

import (
	"messengerServer/internal/services/authorization/api/handlers"
	"messengerServer/internal/services/authorization/config"

	"github.com/gofiber/fiber/v2"
)

func WebStart(conf config.AuthServiceConfig) {
	app := fiber.New()
	app.Get("/authorize", handlers.AuthorizationHandler(conf))
	app.Get("/token", handlers.CheckTokenHandler(conf))
	app.Get("/register", handlers.RegistrationHandler(conf))
	app.Get("/refresh", handlers.RefreshHandler(conf))
	app.Listen((":" + conf.Port))
}
