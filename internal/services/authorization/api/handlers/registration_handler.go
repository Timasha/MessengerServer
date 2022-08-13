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

func RegistrationHandler(conf config.AuthServiceConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var resp responses.RegistrationResponse

		c.Response().Header.SetContentType("application/json")

		if string(c.Request().Header.ContentType()) != "application/json" {
			resp.Err = "wrong content type"
			webUtils.WriteResponse(resp, 406, c)
			return nil
		}

		var (
			req        requests.RegistrationRequest
			unmarshErr error = json.Unmarshal(c.Body(), &req)
		)

		if unmarshErr != nil {
			resp.Err = "cannot unmarshal given JSON: " + unmarshErr.Error()
			webUtils.WriteResponse(resp, 400, c)
			return nil
		} else if len([]byte(req.Login)) < 5 || len([]byte(req.Password)) < 9 {
			resp.Err = "login length lower than 5 or password length lower than 9"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		_, getUserErr := database.GetUser(req.Login)

		if getUserErr != pgx.ErrNoRows && getUserErr != nil {
			resp.Err = "internal database error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		} else if getUserErr == nil {
			resp.Err = "user already exists"
			webUtils.WriteResponse(resp, 400, c)
			return nil
		}

		access, accessGenErr := token.GenerateAccessToken(req.Login, conf.JwtKey, time.Duration(conf.JWTLifetime))
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

		_, createUserErr := database.CreateUser(database.User{Login: req.Login, Password: req.Password, RefreshBodies: []string{refresh[4 : len(refresh)-6]}})

		if createUserErr != nil {
			resp.Err = "create user internal error"
			webUtils.WriteResponse(resp, 500, c)
			return nil
		}

		resp.AccessToken = access
		resp.RefreshToken = refresh
		webUtils.WriteResponse(resp, 201, c)
		return nil
	}
}
