package api

import (
	"errors"
	"log"
	"messengerServer/internal/services/authorization/api/handlers"
	"messengerServer/internal/services/authorization/config"
	"os"

	"github.com/gofiber/fiber/v2"
)

func TlsStart(app *fiber.App, port string) error {
	path := os.Getenv("AUTHSERVICE_PATH_TLS")
	if path == "" {
		return errors.New("tls files path enviroment variable empty or not exists")
	}
	return app.ListenTLS((":" + port), (path + "domain.crt"), (path + "domain.key"))
}

func WebStart(conf config.AuthServiceConfig) {
	app := fiber.New()
	app.Get("/authorize", handlers.AuthorizationHandler(conf))
	app.Get("/token", handlers.CheckTokenHandler(conf))
	app.Get("/register", handlers.RegistrationHandler(conf))
	app.Get("/refresh", handlers.RefreshHandler(conf))
	app.Get("/", func(c *fiber.Ctx) error {
		c.WriteString("hello world")
		return nil
	})
	log.Fatalln(TlsStart(app, conf.Port))

}
