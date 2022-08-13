package handlers

import (
	"messengerServer/internal/api_objects/authorization/responses"
	"messengerServer/internal/services/authorization/config"
	"messengerServer/internal/services/authorization/token"
	"messengerServer/internal/services/authorization/webUtils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CheckTokenHandler(conf config.AuthServiceConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		var resp responses.CheckTokenResponse

		// Установка типа ответа
		c.Response().Header.SetContentType("application/json")

		if string(c.Request().Header.ContentType()) != "application/json" {
			resp.Err = "wrong content type"
			webUtils.WriteResponse(resp, 406, c)
			return nil
		}

		var (
			authHeader string = c.GetReqHeaders()["Authorization"]
		)
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 && parts[0] != "Bearer" {
			resp.Err = "wrong header pattern"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		rawToken := parts[1]

		login, parseErr := token.ParseAccessToken(rawToken, conf.JwtKey)

		if parseErr != nil {
			resp.Err = parseErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		if len([]byte(login)) < 5 {
			resp.Err = "login length lower than 5"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		resp.Login = login
		webUtils.WriteResponse(resp, 200, c)
		return nil
	}
}
