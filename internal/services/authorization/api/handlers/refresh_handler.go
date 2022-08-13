package handlers

import (
	"encoding/json"
	"messengerServer/internal/api_objects/authorization/requests"
	"messengerServer/internal/api_objects/authorization/responses"
	"messengerServer/internal/services/authorization/config"
	"messengerServer/internal/services/authorization/database"
	"messengerServer/internal/services/authorization/token"
	"messengerServer/internal/services/authorization/webUtils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

func RefreshHandler(conf config.AuthServiceConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		var resp responses.RefreshResponse

		// Установка типа ответа
		c.Response().Header.SetContentType("application/json")

		if string(c.Request().Header.ContentType()) != "application/json" {
			resp.Err = "wrong content type"
			webUtils.WriteResponse(resp, 406, c)
			return nil
		}

		// Подготовка данных запроса для обработки
		var (
			req        requests.RefreshRequest
			unmarshErr error = json.Unmarshal(c.Body(), &req)
		)

		if unmarshErr != nil {
			resp.Err = "cannot unmarshal given JSON: " + unmarshErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		if req.RefreshToken[len(req.RefreshToken)-6:] != req.AccessToken[len(req.AccessToken)-6:] {
			resp.Err = "tokens not connected"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		tokenLogin, parseErr := token.ParseAccessToken(req.AccessToken, conf.JwtKey)
		if parseErr != nil && parseErr.Error() != "token is expired" {
			resp.Err = "access is not valid: " + parseErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		dbUser, getUserErr := database.GetUser(tokenLogin)
		if getUserErr == pgx.ErrNoRows {
			resp.Err = "user " + tokenLogin + " not exist"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if getUserErr != nil {
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		refreshIndex := -1
		for i, refreshBody := range dbUser.RefreshBodies {
			if refreshBody == req.RefreshToken[4:len(req.RefreshToken)-6] {
				refreshIndex = i
			}
		}
		if refreshIndex < 0 {
			resp.Err = "refresh not valid"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		access, accessGenErr := token.GenerateAccessToken(dbUser.Login, conf.JwtKey, time.Duration(conf.JWTLifetime))
		if accessGenErr != nil {
			resp.Err = "access token generation internal error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		refresh, refreshGenErr := token.GenerateRefreshToken(access, conf.RefreshLifetime)
		if refreshGenErr != nil {

			resp.Err = "refresh token generation internal error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		if refreshIndex == (len(dbUser.RefreshBodies) - 1) {
			_, updateErr := database.UpdateUser(dbUser.Login, database.User{RefreshBodies: append(dbUser.RefreshBodies[:refreshIndex], refresh[4:len(refresh)-6])})
			if updateErr != nil {
				resp.Err = "internal database error"
				webUtils.WriteResponse(resp, 500, c)
				return nil
			}
		} else {
			_, updateErr := database.UpdateUser(dbUser.Login, database.User{RefreshBodies: append(dbUser.RefreshBodies[:refreshIndex], append(dbUser.RefreshBodies[refreshIndex+1:], refresh[4:len(refresh)-6])...)})
			if updateErr != nil {
				resp.Err = "internal database error"
				webUtils.WriteResponse(resp, 500, c)
				return nil
			}
		}
		resp.AccessToken = access
		resp.RefreshToken = refresh
		return nil
	}
}
