package handlers

import (
	"encoding/json"
	"logger/internal/authService/config"
	"logger/internal/authService/database"
	"logger/internal/authService/responses"
	"logger/internal/authService/token"
	"logger/internal/authService/webUtils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

func AuthorizeHandler(conf config.AuthServiceConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Объект ответа для сериализации в json
		resp := responses.AuthorizeResponce{
			Token: "",
			Err:   "",
		}

		// Проверка MIME типа
		if string(c.Request().Header.ContentType()) != "application/json" {
			resp.Err = "wrong MIME type"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		// Результаты запроса пользователя из базы данных
		var (
			dbUser     database.User
			getUserErr error
		)

		// Подготовка данных запроса для обработки
		var (
			body       []byte        = c.Body()
			user       database.User = database.User{}
			unmarshErr error         = json.Unmarshal(body, &user)
		)

		// Установка типа ответа
		c.Response().Header.SetContentType("application/json")

		// Обработка сырых данных запроса
		if unmarshErr != nil {
			resp.Err = "cannot unmarshal given JSON"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if user.Login == "" || user.Password == "" {
			resp.Err = "empty login or password"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		dbUser, getUserErr = database.GetUser(user.Login)

		if getUserErr == pgx.ErrNoRows { // Проверка на наличие пользователя
			resp.Err = "user not found"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if getUserErr != nil { // Проверка на остальные ошибки
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		} else if dbUser.Password != user.Password { // Если пользователь найден, проверка на правильность пароля
			resp.Err = "wrong password"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		accessToken, tokenGenErr := token.GenerateAccessJWT(user.Login, conf.JwtKey, time.Duration(conf.JWTLifetime))
		if tokenGenErr != nil {
			resp.Err = "cannot generate authorization token"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		resp.Token = accessToken

		webUtils.WriteResponse(resp, 200, c)

		return nil
	}
}
