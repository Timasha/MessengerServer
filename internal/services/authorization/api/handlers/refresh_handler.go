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
		if getUserErr == database.ErrNoRows {
			resp.Err = "user " + tokenLogin + " not exist"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if getUserErr != nil {
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}

		refreshIndex, validErr := token.ValidRefreshToken(req.RefreshToken, req.AccessToken, dbUser.RefreshBodies)

		if validErr.Error() == "refresh is expired" {
			_, updateErr := database.UpdateUser(dbUser.Login, database.User{RefreshBodies: append(dbUser.RefreshBodies[:refreshIndex], dbUser.RefreshBodies[:refreshIndex+1]...)})
			if updateErr != nil {
				resp.Err = "internal database error"
				webUtils.WriteResponse(resp, 500, c)
				return nil
			}
			resp.Err = validErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}
		if validErr != nil {
			resp.Err = "refresh not valid: " + validErr.Error()
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
		dbUser.RefreshBodies[refreshIndex] = refresh[4 : len(refresh)-6]
		_, updateErr := database.UpdateUser(dbUser.Login, database.User{RefreshBodies: dbUser.RefreshBodies})
		if updateErr != nil {
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}
		resp.AccessToken = access
		resp.RefreshToken = refresh
		return nil
	}
}
