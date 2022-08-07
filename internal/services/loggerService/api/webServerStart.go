package api

import (
	"messengerServer/internal/services/loggerService/config"

	"github.com/gofiber/fiber/v2"
)

func WebStart(conf config.LoggerServiceConfig) {
	app := fiber.New()
	app.Listen((":" + conf.Port))
}
