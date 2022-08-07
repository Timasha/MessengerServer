package handlers

import (
	"encoding/json"
	"messengerServer/internal/services/authService/config"
	"messengerServer/internal/services/authService/database"
	"messengerServer/internal/services/authService/token"
	"messengerServer/internal/services/authService/webUtils"
	"messengerServer/internal/servicesApi/authService/requests"
	"messengerServer/internal/servicesApi/authService/responses"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

func AuthorizationHandler(conf config.AuthServiceConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Объект ответа для сериализации в json
		var resp responses.AuthorizationResponse

		// Установка типа ответа
		c.Response().Header.SetContentType("application/json")

		// Проверка MIME типа
		if string(c.Request().Header.ContentType()) != "application/json" {
			resp.Err = "wrong content type"
			webUtils.WriteResponse(resp, 406, c)
			return nil
		}

		// Результаты запроса пользователя из базы данных
		var (
			dbUser     database.User
			getUserErr error
		)

		// Подготовка данных запроса для обработки
		var (
			req        requests.AuthorizeRequest
			unmarshErr error = json.Unmarshal(c.Body(), &req)
		)

		// Обработка сырых данных запроса
		if unmarshErr != nil {
			resp.Err = "cannot unmarshal given JSON: " + unmarshErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if len([]byte(req.Login)) < 5 || len([]byte(req.Password)) < 9 {
			resp.Err = "login length lower than 5 or password length lower than 9"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		dbUser, getUserErr = database.GetUser(req.Login)

		if getUserErr == pgx.ErrNoRows { // Проверка на наличие пользователя
			resp.Err = "user not found"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if getUserErr != nil { // Проверка на остальные ошибки
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		} else if dbUser.Password != req.Password { // Если пользователь найден, проверка на правильность пароля
			resp.Err = "wrong password"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		accessToken, tokenGenErr := token.GenerateAccessToken(req.Login, conf.JwtKey, time.Duration(conf.JWTLifetime))
		if tokenGenErr != nil {
			resp.Err = "access token generation internal error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		refresh, refreshGenErr := token.GenerateRefreshToken(accessToken, conf.RefreshLifetime)
		if refreshGenErr != nil {
			resp.Err = "refresh token generation internal error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		_, addRefreshErr := database.AddRefresh(req.Login, refresh[4:len(refresh)-6])
		if addRefreshErr != nil {
			resp.Err = "cannot add refresh to database"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}

		resp.AccessToken = accessToken
		resp.RefreshToken = refresh
		webUtils.WriteResponse(resp, 200, c)

		return nil
	}
}
