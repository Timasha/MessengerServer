package webUtils

import (
	"encoding/json"
	"errors"
	"logger/internal/authService/responses"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// Функция для отправки ответа на запрос в зависимости от его типа.
func WriteResponse(resp interface{}, respCode int, c *fiber.Ctx) error {
	switch resp.(type) {
	case responses.AuthorizeResponce:
		{
			rawResp, _ := json.Marshal(resp)
			_, err := c.Status(respCode).Write(rawResp)
			return err
		}
	default:
		return errors.New("wrong response type: " + reflect.TypeOf(resp).String())
	}
}
